import { Button } from "@/components/ui/button";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useState, useEffect } from "react";

export const Route = createFileRoute("/queue/$queueId")({
  component: RouteComponent,
  loader: async ({ params }) => {
    const resp = await fetch(`/api/queue/${params.queueId}`);
    if (resp.status === 404) {
      throw new Error(`queue ${params.queueId} not found`);
    }

    const data = await resp.json();

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

interface QueueStatus {
  chair_list: string[];
  queue_list: number[];
  ready: boolean;
  checked_in: boolean;
}

function RouteComponent() {
  const params = Route.useParams();
  const queueNumber = parseInt(params.queueId);
  const { queueInfo } = Route.useLoaderData();

  const [queueStatus, setQueueStatus] = useState<QueueStatus | undefined>();

  useEffect(() => {
    const url = `/api/queue/${queueNumber}/stream-status`;
    const evSource = new EventSource(url);
    evSource.addEventListener("message", (ev) => {
      setQueueStatus(JSON.parse(ev.data));
    });
  }, []);

  const checkIn = async () => {
    console.log("checkin");
    // const resp = await fetch(`/api/queue/2/check-in`, {
    const resp = await fetch(`/api/queue/${params.queueId}/check-in`, {
      method: "POST",
    });
  };

  return (
    <div>
      <p>Queue number {queueNumber}</p>
      <p>Name: {queueInfo.name}</p>
      <p>Number: {queueInfo.number}</p>
      <p>Queue status: {JSON.stringify(queueStatus)}</p>
      {queueStatus?.ready ? (
        <>
          <Button onClick={checkIn}>Check In</Button>
        </>
      ) : (
        <p>Waiting</p>
      )}
    </div>
  );
}
