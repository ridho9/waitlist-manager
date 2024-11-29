import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/queue/$queueId")({
  component: RouteComponent,
});

function RouteComponent() {
  const params = Route.useParams();
  const queueNumber = parseInt(params.queueId);
  return <div>Queue number {queueNumber}</div>;
}
