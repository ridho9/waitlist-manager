import { Button } from "@/components/ui/button";
import { PlaceStatus, listenPlaceStatusChange } from "@/lib/placeStatus";
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
    useEffect(() => {});
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

function RouteComponent() {
  const params = Route.useParams();
  const queueNumber = parseInt(params.queueId);
  const { queueInfo } = Route.useLoaderData();

  const [queueStatus, setQueueStatus] = useState();

  useEffect(() => {
    const url = `/api/queue/${queueNumber}/stream-status`;
    const evSource = new EventSource(url);
    evSource.addEventListener("message", (ev) => {
      setQueueStatus(JSON.parse(ev.data));
    });
  }, []);

  return (
    <div>
      <p>Queue number {queueNumber}</p>
      <p>Name: {queueInfo.name}</p>
      <p>Number: {queueInfo.number}</p>
      <p>Queue status: {JSON.stringify(queueStatus)}</p>
    </div>
  );
}
