import { PlaceStatusData } from "@/components/placeStatus";
import { Button } from "@/components/ui/button";
import { createFileRoute, Link } from "@tanstack/react-router";
import { Users, Clock, User } from "lucide-react";
import { useState, useEffect } from "react";
import { QRCode } from "react-qrcode-logo";

export const Route = createFileRoute("/queue/$queueId")({
  component: RouteComponent,
  loader: async ({ params }) => {
    const resp = await fetch(`/api/queue/${params.queueId}`);
    if (resp.status === 404) {
      throw new Error(`queue ${params.queueId} not found`);
    }

    const data = (await resp.json()) as QueueInfo;

    return { queueInfo: data };
  },
  errorComponent: ({ error }) => {
    // Render an error message
    return (
      <div>
        <p>{error.message}</p>
        <Button asChild variant="outline">
          <Link to="/book">Go To Booking Page</Link>
        </Button>
      </div>
    );
  },
});

type QueueStatusData = PlaceStatusData & {
  ready: boolean;
  checked_in: boolean;
};

interface QueueInfo {
  name: string;
  number: number;
  checked_in: boolean;
}

function RouteComponent() {
  const params = Route.useParams();
  const { queueInfo } = Route.useLoaderData();

  return (
    <main className="min-h-screen bg-gray-100 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <QueueStatus
          queueId={params.queueId}
          name={queueInfo.name}
          diners={queueInfo.number}
        />
      </div>
    </main>
  );
}

function QueueStatus(params: {
  queueId: string;
  name: string;
  diners: number;
}) {
  const { queueId } = params;
  const [queueStatus, setQueueStatus] = useState<QueueStatusData | undefined>();

  useEffect(() => {
    const url = `/api/queue/${queueId}/stream-status`;
    const evSource = new EventSource(url);
    evSource.addEventListener("message", (ev) => {
      setQueueStatus(JSON.parse(ev.data));
    });
    return () => {
      evSource.close();
    };
  }, []);

  const filledChairs = queueStatus?.chair_list.length;
  const maxChairs = queueStatus?.max_chair;
  let queueAhead = 0;
  for (const q of queueStatus?.queue_list || []) {
    if (`${q}` === queueId) break;
    queueAhead += 1;
  }

  const checkIn = async () => {
    return fetch(`/api/queue/${params.queueId}/check-in`, {
      method: "POST",
    });
  };

  return (
    <>
      <div className="bg-white shadow rounded-lg p-6 text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Queue Status</h1>
        <div
          className="text-4xl font-bold text-primary mb-8"
          aria-live="polite"
        >
          #{queueId}
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div className="bg-gray-50 p-4 rounded-lg">
            <div className="flex items-center justify-center mb-2">
              <Users className="h-6 w-6 text-primary mr-2" />
              <h2 className="text-xl font-semibold text-gray-900">
                Filled Chairs
              </h2>
            </div>
            <p className="text-3xl font-bold text-gray-700">
              {filledChairs}/{maxChairs}
            </p>
          </div>
          <div className="bg-gray-50 p-4 rounded-lg">
            <div className="flex items-center justify-center mb-2">
              <Clock className="h-6 w-6 text-primary mr-2" />
              <h2 className="text-xl font-semibold text-gray-900">
                Queue Ahead
              </h2>
            </div>
            <p className="text-3xl font-bold text-gray-700">
              {queueStatus?.checked_in ? "0" : queueAhead}
            </p>
          </div>
        </div>
        <div className="bg-gray-50 p-4 rounded-lg mb-6">
          <div className="flex items-center justify-center mb-2">
            <User className="h-6 w-6 text-primary mr-2" />
            <h2 className="text-xl font-semibold text-gray-900">
              Your Reservation
            </h2>
          </div>
          <p className="text-xl text-gray-700">
            <span className="font-semibold">{params.name}</span> -{" "}
            {params.diners} {params.diners === 1 ? "person" : "people"}
          </p>
        </div>
        {queueStatus ? (
          <div className="bg-gray-50 p-4 rounded-lg mb-6">
            <WaitComponent
              queueId={queueId}
              queueStatus={queueStatus as QueueStatusData}
              name={params.name}
              diners={params.diners}
              onCheckIn={checkIn}
            />
          </div>
        ) : undefined}
        <div className="bg-gray-50 p-4 rounded-lg mb-6">
          <div className="flex flex-col items-center justify-center mb-2">
            <p>Share This Queue</p>
            <p>Open {window.location.href} or scan QR below</p>
            <QRCode value={window.location.href} />
          </div>
        </div>
      </div>
    </>
  );
}

const WaitComponent = (params: {
  queueId: string;
  queueStatus: QueueStatusData;
  name: string;
  diners: number;
  onCheckIn: () => Promise<any>;
}) => {
  const { queueStatus, name, diners, queueId } = params;
  const isYourTurn = `${queueStatus.queue_list[0]}` === queueId;
  const enoughChairs =
    queueStatus.max_chair - queueStatus.chair_list.length >= diners;
  const [buttonEnabled, setButtonEnabled] = useState(true);

  if (params.queueStatus.checked_in) {
    return (
      <>
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Status</h2>
        <p className="text-lg font-medium">Checked In, Enjoy Your Meal</p>
      </>
    );
  }

  return (
    <>
      <h2 className="text-xl font-semibold text-gray-900 mb-4">Status</h2>
      {!isYourTurn && (
        <p className="text-lg text-yellow-600 font-medium">
          Waiting for your turn
        </p>
      )}
      {isYourTurn && !enoughChairs && (
        <p className="text-lg text-orange-600 font-medium">
          Waiting for enough chairs
        </p>
      )}
      {isYourTurn && enoughChairs && queueStatus.ready && (
        <Button
          className="w-full text-lg"
          size="lg"
          disabled={!buttonEnabled}
          onClick={async () => {
            setButtonEnabled(false);
            await params.onCheckIn();
          }}
        >
          Check In
        </Button>
      )}
    </>
  );
};
