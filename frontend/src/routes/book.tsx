import * as React from "react";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/book")({
  component: BookComponent,
});

interface PlaceStatus {
  chair_list: string[];
}

interface QueueResp {
  queue_number: number;
  queue_list: number[];
}

function BookComponent() {
  const [assignedQueue, setAssignedQueue] = useState<QueueResp | undefined>();
  const navigate = useNavigate({ from: "/book" });

  const registerForm = (
    <PartyRegisterForm
      onSubmit={async (data) => {
        console.log(data);
        const body = {
          party_name: data.partyName,
          party_number: data.partyNumber,
        };
        const resp = await fetch(`/api/queue`, {
          method: "POST",
          body: JSON.stringify(body),
        });
        const respBody = await resp.json();
        setAssignedQueue(respBody);
        console.log("start timeout");
        await new Promise((res, rej) => setTimeout(res, 3000));
        navigate({
          to: "/queue/$queueId",
          params: { queueId: `${respBody.queue_number}` },
        });
        console.log("redirecting");
      }}
    />
  );

  const confimAdded = (
    <>
      <p>Successfully added to the queue!</p>
      <p>
        Queue Number:{" "}
        <span className="font-bold">{assignedQueue?.queue_number}</span>
      </p>
      <p>Redirecting to waiting page...</p>
    </>
  );

  return (
    <div className="p-2">
      <h1 className="font-bold text-xl">Welcome to Restaurant</h1>
      <div className="py-2">
        <PlaceStatus />
      </div>

      <div className="max-w-sm my-2">
        {assignedQueue ? confimAdded : registerForm}
      </div>
    </div>
  );
}

function PlaceStatus() {
  const [placeStatus, setPlaceStatus] = useState<PlaceStatus>();

  useEffect(() => {
    const url = `/api/stream-place-status`;
    const evSource = new EventSource(url);
    evSource.addEventListener("message", (ev) => {
      setPlaceStatus(JSON.parse(ev.data));
    });
  }, []);

  return (
    <>
      <p className="font-bold text-lg">Current Restaurant Status</p>
      <div>{JSON.stringify(placeStatus)}</div>
    </>
  );
}

interface PartyRegisterData {
  partyName: string;
  partyNumber: number;
}

function PartyRegisterForm({
  onSubmit,
}: {
  onSubmit: (data: PartyRegisterData) => Promise<void>;
}) {
  const [partyName, setPartyName] = useState("");
  const [partyNumber, setPartyNumber] = useState(1);
  const [buttonEnable, setButtonEnable] = useState(true);

  return (
    <form
      onSubmit={async (e) => {
        e.preventDefault();
        setButtonEnable(false);
        await onSubmit({ partyName, partyNumber });
        setButtonEnable(true);
      }}
    >
      <p className="font-bold text-lg">Register your party</p>
      <Input
        name="name"
        type="text"
        placeholder="Name"
        value={partyName}
        onChange={(ch) => setPartyName(ch.currentTarget.value)}
        required
        className="my-1"
      />
      <Input
        name="number"
        type="number"
        value={partyNumber}
        min={1}
        max={10}
        placeholder="1"
        className="my-1"
        onChange={(e) => setPartyNumber(e.currentTarget.valueAsNumber)}
      />
      <p>{buttonEnable}</p>
      <Button type="submit" disabled={!buttonEnable}>
        {buttonEnable ? "Register" : "Please Wait"}
      </Button>
    </form>
  );
}