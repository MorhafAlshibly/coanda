---
id: schema
slug: /
title: coanda
sidebar_position: 1
hide_table_of_contents: false
pagination_next: null
pagination_prev: null
sidebar_class_name: navbar__toggle
---
coanda is a modern API for the modern game developer. It is a GraphQL API that provides a simple and efficient way to access the data you need to build your game. It is designed to be easy to use, fast, and reliable. It is built on top of the latest technologies and is constantly updated to provide the best experience for game developers.

## 🗓️ Events

The **Event microservice** is responsible for managing single-use competitive events. It handles the full lifecycle of an event, from creation and round setup, to user participation, scoring, and final result dispatching.

---

### 🚀 Purpose

This microservice enables time-bound game events where users compete across multiple rounds. It supports user participation, scoring logic, leaderboard aggregation, and final event submission to external systems for further processing.

---

### 🧱 Core Concepts

#### **Event**

- A **single-use competition instance**.
- Defined by a **start time**.
- Contains one or more **Rounds**.
- Ends after the last round’s `ended_at` timestamp, however the event still will exist in the table.
- Sends final results to a configured **third-party target** (e.g., game backend).

#### **Round**

- Represents a **phase or stage** within an event.
- Defined by a unique `ended_at` time (no two rounds can end at the same timestamp).
- Includes a **scoring array**, which specifies the points assigned to each placement (e.g., `[10, 5, 3]` → 1st = 10 pts, 2nd = 5 pts, etc.).
- Every user who finishes a placement past the length of the scoring array will default to a score of 0.

#### **User Participation**

- Users are added to an event when the AddEventResult endpoint is called with their result for a specific round.
- They can then submit their results for the rest of the rounds they are a part of.
- Submissions are used to build **leaderboards**, either per round, which sorts results in ascending order with the lowest result being 1st place.
- Or a event wide leaderboard, which sorts scores in descending order, with the highest score being 1st place.
- All rankings are done via a dense rank function, meaning we may have multiple 1st place users, or multiple 5th place users, if scores/results are the same.

## 📦 Items

The **Item microservice** provides a lightweight, general-purpose storage solution for structured or unstructured data. It allows other services to create, store, query, and automatically expire arbitrary data objects, known as **items**.

---

### 🚀 Purpose

This service acts as a flexible key-value store where each item is uniquely identified by a combination of an `id` and a `type`. It’s useful for storing metadata, temporary state, session data, or any other JSON-based payloads that need to be queried or expire after a certain time.

---

### 🧱 Core Concepts

#### **Item**

- Uniquely identified by a combination of:
  - `id` (string)
  - `type` (string)
- Stores an arbitrary `data` field as raw JSON.
- Includes an optional `expires_at` timestamp to define when the item should be automatically removed from the database.

---

### 🛠️ Features

- **Flexible Schema:** The `data` field accepts any valid JSON object, structured or unstructured.
- **Type Filtering:** Items can be queried by their `type`, allowing consumers to narrow down to specific categories of items.
- **Auto Expiry:** Items with an `expires_at` value are automatically excluded after expiration.

## 🎯 Matchmaking

The **Matchmaking microservice** is responsible for pairing users into matches based on ELO ratings, time-based queue dynamics, and arena rules. It handles users, tickets, arenas, and match creation, and tries to be a fair, scalable, and flexible game session manager.

---

### 🧱 Core Concepts

#### **Arena**

- Represents a **game mode or venue** users can queue into.
- Defines matchmaking rules:
  - `max_players` – max number of players per match.
  - `min_players` – minimum number required to **start** a match.
  - `max_players_per_ticket` – max number of players allowed in a single ticket.

#### **User**

- A player with a unique identifier and an **ELO rating**.
- Can only exist in **one ticket at a time**.
- Must be explicitly deleted when no longer queued or matched.

#### **Ticket**

- A matchmaking request from a user or group of users.
- Tied to a single **arena** or multiple compatible arenas.
- Represents a **unit of matchmaking**.
- An evolving **ELO window** is tracked based on the creation time of a ticket, allowing for broader match acceptance.
- Only deleted **manually** - you must delete it if you don't want it to be valid anymore
  - If it is matched you must delete the match its a part of to delete it, this is something to take note of when clients suddenly "crash", yet tickets stay in the queue.

#### **Match**

- Represents a **group of tickets** selected to play together.
- Tied to an **arena**, inherits its capacity rules.
- Can be assigned a **private server ID** (set only once) to sync with a game backend for deployment.
- Will only start once it satisfies `min_capacity` and a private server ID is set.
- Once started:
  - It **locks** and no longer accepts tickets.
- Must be manually deleted; doing so cascades and removes its associated tickets and users.

---

### ⚙️ Matchmaking Flow

1. A **user** or **group of users** is created with their elos.
2. A **ticket** is created for a **group of users** and a **group of arenas**.
3. **Ticket waits** in the queue, attempting to find matches with **similar average ELO**.
4. **ELO window expands** over time to include a broader pool.
5. If no match is found:
   - A **match is created** (with that ticket).
   - It **waits** for viable additional tickets to arrive.
6. The game backend can then assign a private server ID to the match.
7. Once the match is started, no more tickets can be assigned to the match.
8. When the match ends, you can set a end time and/or delete the match

---

### 🔁 Manual Deletion and Cascade Rules

- **Users**, **tickets**, and **matches** are not auto-deleted.
- The game backend is responsible for cleanup.
- **Cascading deletes**:
  - Deleting a **match** → deletes all associated **tickets** and **users**.
  - Deleting a **ticket** → deletes the associated **user**.
  - A user cannot be deleted if it is assigned to a ticket, ticket must be deleted.
    - Likewise with tickets and matches.

### 🪦 Abandoned Tickets

There may exist scenarios where you create a ticket for a user or group of users, then fail to keep track of the ticket if for example your game backend or server shuts down. These tickets will still stay in the system. They may either:

* Join an existing match
* Create its own match if it can't find an existing match in enough time.

#### How do you handle it?

* Before starting any match ensure every ticket in the match is **"alive"**, meaning the users are real users and not abandoned.
  * If a ticket exists that isn't **alive**, then delete the match and requeue the other tickets.
* If the ticket **creates its own match**, you likely won't be tracking it.
  * You may begin tracking it once an **alive** ticket joins that match, allowing you to **delete the match** as mentioned previously.
  * Or if the users of the match attempt to create a **new** ticket, you notice that they're currently in a ticket/match, and can then **delete the match**.


## 🏅 Records

The **Record microservice** tracks user performance in specific game modes and provides a way to query ranked leaderboards. Each record represents a user's best result in a given game mode and is ranked against others.

---

### 🚀 Purpose

This service is designed to store and manage users' **best performances** in competitive or trackable game modes (e.g. speedruns, high scores, time trials). It also supports leaderboard generation and efficient updates to users' existing records.

---

### 🧱 Core Concepts

#### **Record**

- Represents a **user's best score/time/performance** in a specific game mode.
- Identified by:
  - `id` – unique identifier for the record.
  - `user_id` – the player this record belongs to.
  - `name` – the game mode or category the record is for.
  - `record` – the actual performance value (time).
  - `ranking` – where this record sits relative to others in the same `name` (ranked densely in ascending order).
- A user can only have **one record per `name`**, ensuring that only their best is stored.

---

### 📊 Leaderboards

- Records can be **queried by `name`**, returning a **leaderboard** sorted by the `record` value in ascending order.
- Rankings are automatically recalculated based on performance comparisons across users in the same game mode.

## ✅ Tasks

The **Task microservice** is a lightweight, general-purpose service for tracking user or system-defined tasks with support for completion status. It extends the functionality of the `Item` microservice by adding a `completed` field and completion logic.

---

### 🚀 Purpose

This service is ideal for managing **checklist-like functionality**, whether for user missions, background jobs, daily goals, or any task-based logic within a game or system. It allows for the creation, filtering, and completion of tasks.

---

### 🧱 Core Concepts

#### **Task**

- Identified by a unique combination of `id`, and `type`.
- Stores:
  - `data` – arbitrary JSON payload (structured or unstructured).
  - `expires_at` – optional timestamp after which the task is no longer valid.
  - `completed` – boolean flag (default `false`) indicating task completion.

---

### 🛠️ Features

- **Flexible Structure:** Like the `Item` service, stores arbitrary JSON in the `data` field.
- **Completion Support:**
  - Tasks can be **marked as completed** via the API.
  - Can be **filtered** by `completed=true|false` when querying.
- **Type Filtering:** Query by `type` to group or subset tasks logically.
- **Expiry:** Tasks with `expires_at` are ignored/removed after expiration time.

## 🧑‍🤝‍🧑 Teams

The **Team microservice** manages groups of users (teams) and their memberships, along with scoring and ranking across all teams.

---

### 🚀 Purpose

This service enables creation and management of **teams** with membership limits, tracks team scores, and provides ranked leaderboards based on team performance.

---

### 🧱 Core Concepts

#### **Team**

- Identified by a unique `name`.
- Teams have a limit set in the application of how many members max can be in a team.
- Has:
  - A `score` representing team performance.
  - A `ranking` based on the team’s score in ascending order using **dense ranking**.

#### **Team Member**

- Represents a user who is a member of exactly **one** team at a time.
- Users cannot belong to multiple teams simultaneously.
- Membership enforces the team's member limit.

---

### 🧮 Scoring & Ranking

- Teams are ranked by their `score` in **ascending order** (lower score = better rank).
- **Dense ranking** is used, so ties share the same rank, and ranks increment without gaps.

---

### 👥 Membership Rules

- Team members are restricted to a **single team**.
- Adding a member to a team checks that the team's current member count is below the `limit`.
- Membership updates (add/remove) are managed via API calls.

## 🏆 Tournaments

The **Tournament microservice** manages recurring competitive events where users participate and accumulate scores over defined intervals.

---

### 🚀 Purpose

This service tracks user scores within tournaments that repeat on a set schedule, such as daily, weekly, or monthly cycles. It automatically resets results at the end of each interval and notifies third-party systems for further processing.

---

### 🧱 Core Concepts

#### **Tournament User**

- Represents a user's participation and score in a specific tournament.
- Contains:
  - `id` – the unique ID of the tournament user.
  - `tournament` – string identifier of the tournament.
  - `user_id` – unique ID of the user (for the game backend).
  - `score` – the user’s current score in the tournament.
  - `interval` – defines the tournament's reset frequency.

#### **Tournament Interval**

- Defines how often the tournament resets and results are wiped:
  - `DAILY`
  - `WEEKLY`
  - `MONTHLY`
  - `UNLIMITED` (no automatic resets)

---

### 🔄 Lifecycle & Resetting

- Tournaments run continuously but are segmented into intervals based on the `interval` setting.
- At the end of each interval (e.g., day, week, month):
  - Aggregated results are sent to a configured third-party source (e.g., game backend) for processing.
  - A new tournament interval begins automatically.
  - However, the actual tournament data may not be deleted yet, it could still be queryable by ID, however it is hidden when querying by `name`, `user_id`, `interval`.

---

### ℹ️ Implementation Details

- Only the **Tournament User** table exists.
- The active tournament for a user is inferred automatically from the user’s creation timestamp and the tournament’s interval.
- This simplifies storage and querying by avoiding explicit tournament instances.

## 🪝Webhook

There is a webhook function that allows you to use the API as a proxy to fulfill webhook requests, be careful however with user submitted data, as it can lead to malicious use of the API.
