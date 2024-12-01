import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Utensils } from "lucide-react";
import { Label } from "@/components/ui/label";
import { PlaceStatus } from "@/components/placeStatus";

export const Route = createFileRoute("/book")({
  component: BookPage,
});

interface QueueResp {
  queue_number: number;
  queue_list: number[];
}

function BookPage() {
  return (
    <>
      <main className="min-h-screen bg-gray-100 py-8 px-4 sm:px-6 lg:px-g">
        <div className="max-w-3xl mx-auto space-y-8">
          <Header />
          <PlaceStatus />
          <QueueForm />
        </div>
      </main>
    </>
  );
}

function Header() {
  return (
    <header className="text-center">
      <div className="flex items-center justify-center">
        <Utensils className="h-8 w-8 text-primary mr-2" />
        <h1 className="text-3xl font-bold text-gray-900">
          Welcome to Our Restaurant
        </h1>
      </div>
    </header>
  );
}

function QueueForm() {
  const navigate = useNavigate({ from: "/book" });

  const [name, setName] = useState("");
  const [diners, setDiners] = useState("1");
  const [buttonEnable, setButtonEnable] = useState(true);
  const [assignedQueue, setAssignedQueue] = useState<QueueResp | undefined>();

  const handleSubmit = async (ev: React.FormEvent) => {
    ev.preventDefault();
    setButtonEnable(false);

    const body = {
      party_name: name,
      party_number: parseInt(diners),
    };
    const resp = await fetch(`/api/queue`, {
      method: "POST",
      body: JSON.stringify(body),
    });
    const respBody = await resp.json();
    setAssignedQueue(respBody);
    await new Promise((res, rej) => setTimeout(res, 1000));
    await navigate({
      to: "/queue/$queueId",
      params: { queueId: `${respBody.queue_number}` },
    });

    setButtonEnable(true);
  };

  if (assignedQueue) {
    return (
      <div className="bg-white shadow rounded-lg p-6 text-center">
        <h2 className="text-2xl font-semibold text-gray-900 mb-4">
          Queue Joined Successfully!
        </h2>
        <p className="text-xl text-gray-700 mb-4">
          Your Queue Number is{" "}
          <span className="font-bold text-primary">
            {assignedQueue.queue_number}
          </span>
        </p>
        <p className="text-gray-600">Redirecting...</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h2 className="text-xl font-semibold text-gray-900 mb-4">
        Join the Queue
      </h2>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <Label htmlFor="name">Name</Label>
          <Input
            type="text"
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            className="mt-1"
          />
        </div>
        <div>
          <Label htmlFor="diners">Number of Diners</Label>
          <Input
            type="number"
            id="diners"
            value={diners}
            onChange={(e) => setDiners(e.target.value)}
            required
            min="1"
            max="10"
            className="mt-1"
          />
        </div>
        <Button type="submit" className="w-full" disabled={!buttonEnable}>
          Join Queue
        </Button>
      </form>
    </div>
  );
}
