export interface PlaceStatus {
  chair_list: string[];
  queue_list: string[];
}

export function listenPlaceStatusChange(
  onChange: (status: PlaceStatus) => void
) {
  const url = `/api/stream-place-status`;
  const evSource = new EventSource(url);
  evSource.addEventListener("message", (ev) => {
    onChange(JSON.parse(ev.data));
  });
}
