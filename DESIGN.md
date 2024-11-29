# Flow

## 1. User opens main app

Page: `/book`. This page will be the main landing page that the user will open.

Further considerations: if later app used by multiple restaurant, the endpoint can be `/r/:id/book` where `:id` is a unique restaurant ID.

### Connections:

1. `GET /api/stream-place-status`

This is an SSE endpoint used to stream the data to the `/book` page so that the user can see the restaurant status in real time, mainly the amount of available chairs, and the existing queue.

`stream-` prefix is used on the route to signify that it's a SSE endpoint.

`data` payload
```ts
{
    "total_seats": number;
    "available_seats": number;
    "queue_status": {
        "amount_waiting": number;
    };
}
```

2. `POST /api/queue`

```ts
{
    name: string;
    party_size: number;
}
```

This api is used to add new party to the queue. Response

```ts
{
    queue_id: string;
}
```

`queue_id` is later used to track the queue.

After this user is redirected to the queue waiting page.

## 2. User opens waiting page

`/queue/:queue_id`

This page shows:
1. Name of party
2. Amount of people
3. queue ahead

Status streams from SSE endpoint `/api/stream-queue-status?queue-id={queue_id}` (afterward referred as queue status endpoint).
This endpoint streams these events:

1. `event: status-update`

data:
```ts
{
    name: string;
    party_size: number;
    party_amount_ahead: number;
}
```
2. `event: party-ready`

This event is sent by the backend when the current party is ready to check-in. After this event is received, the front-end should show the check-in button.

3. `event: party-checked-in`

When one of the devices observing this page press check in, the endpoint will send this event to note that the party is checked-in. This event is also sent once if hte queue waiting page opened after checked in

## 3. Users check in

After user press the check-in button above, they will send a checkin request.

`POST /api/queue/:queue-id/check-in`

This endpoint will check whether the party is actually ready, and seats them.

The `/queue/:queue-id` will show that the party is checked in, and after service is finished, will check-out the party.


# System Design

Our database will be mainly a Redis, with the consideration that this particular type of app is heavily event-driven, and Redis have a good support for simple pubsub system. Redis also have a queue data structure that we can use.

Chair model:

The restauran chair is stored as a list with that is allowed to have up to the max number of chair in restaurant (10 for this case):

The list will be stored in key `chair`. The stored value will be `diner:{queue-id}:{member-count}`.

