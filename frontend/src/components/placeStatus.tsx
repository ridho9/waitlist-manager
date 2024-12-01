import { Users, Clock } from "lucide-react";
import { useState, useEffect } from "react";

export interface PlaceStatusData {
  chair_list: string[];
  queue_list: number[];
  max_chair: number;
}

export function PlaceStatus() {
  const [placeStatus, setPlaceStatus] = useState<PlaceStatusData>();

  useEffect(() => {
    const url = `/api/stream-place-status`;
    const evSource = new EventSource(url);
    evSource.addEventListener("message", (ev) => {
      setPlaceStatus(JSON.parse(ev.data));
    });

    return () => {
      evSource.close();
    };
  }, []);

  return (
    <div className="bg-white shadow rounded-lg p-6">
      <h2 className="text-xl font-semibold text-gray-900 mb-4">
        Current Restaurant Status
      </h2>
      <div className="flex justify-between items-center">
        <div className="flex items-center">
          <Users className="h-6 w-6 text-primary mr-2" />
          <span className="text-lg text-gray-700">
            {placeStatus?.chair_list.length}/{placeStatus?.max_chair} chairs
            filled
          </span>
        </div>
        <div className="flex items-center">
          <Clock className="h-6 w-6 text-primary mr-2" />
          <span className="text-lg text-gray-700">
            {placeStatus?.queue_list.length} in queue
          </span>
        </div>
      </div>
    </div>
  );
}
