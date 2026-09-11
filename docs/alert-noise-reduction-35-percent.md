# How I Cut Alert Noise 35% Without Muting Failures

A practical Grafana alerting teardown for SREs

## TL;DR

The biggest reduction did not come from changing thresholds, adding longer `for` durations, or deleting alert rules. It came from fixing notification semantics.

The alert rules kept evaluating as before. The change was in how active alerts were grouped, routed, and repeated to humans.

The concrete policy moved from alert-name-centric grouping and frequent reminders to service-context grouping with a scoped long repeat interval:

```text
Before
  group_by:         [grafana_folder, alertname]
  repeat_interval:  30m

After, scheduler route only
  group_by:         [grafana_folder, environment, alert_type, application]
  group_wait:        2m
  group_interval:    30m
  repeat_interval:   720h   # 30 days
```

The measured outcome of the tuning cycle was a **35% reduction in human-visible alert noise**.

The important part is what did *not* change: the underlying scheduler failure detection remained active. This was a delivery-layer optimization, not a detection-layer compromise.

---

## The failure mode: correct alerts, bad notification semantics

A noisy alerting system is not always a bad detection system.

In this case, the failure rules were doing their job. The problem was that notification behavior made one operational condition look like many separate incidents.

A typical pattern looked like this:

1. Several scheduler alerts entered `Firing` around the same time.
2. They shared the same application and operational failure context.
3. The notification policy grouped too narrowly around individual alert names.
4. Active alerts were reminded every 30 minutes.
5. Engineers received repeated messages even when there was no new state change or new action to take.

That creates two distinct kinds of noise:

- **fan-out noise**: one underlying problem produces several notifications
- **persistence noise**: an unchanged firing state keeps producing reminders

The fix needs to treat those separately.

---

## First principle: measure notifications, not evaluations

A rule evaluating every 10 minutes is not inherently noisy.

An evaluation is machine work. A notification is human interruption.

That distinction matters because teams often "fix" alert fatigue by weakening detection logic. They increase thresholds, extend lookback windows, or add long pending periods. That can lower page volume, but it can also increase detection latency.

For this tuning cycle, I treated the human-visible notification as the unit of noise.

Conceptually:

```text
notification_noise_reduction =
    1 - (post_change_notifications / pre_change_notifications)
```

The alert generation logic stayed stable while notification policy changed. The observed result was **35% fewer human-visible alert notifications**.

I am intentionally not inventing an absolute before/after count here. I do not have the original auditable raw count export in this public write-up, so a fabricated "1,000 alerts became 650" chart would make the post look more precise while making it less trustworthy.

What I *can* publish are the actual policy values, rule-scale numbers, and validation results used in the change.

---

## The environment

One non-production scheduler environment had **91 execution-failure alert rules**.

Those rules evaluated every **10 minutes** against a **15-minute lookback**, with a simple failure condition equivalent to:

```text
A > 0
```

where `A` represented detected scheduler errors in the lookback window.

The important design choice was to leave this detection layer intact and change only the notification policy.

---

## Change 1: group by operational identity, not alert identity

The original grouping centered on the individual alert name:

```text
[grafana_folder, alertname]
```

That is easy to understand, but it is often the wrong abstraction for an operator.

An on-call engineer usually wants to answer:

- Which application is affected?
- In which environment?
- What kind of failure is this?
- Are these alerts part of the same operational event?

So the grouping key became:

```text
[grafana_folder, environment, alert_type, application]
```

This shifts the grouping model from "which rule fired?" to "which service condition am I responding to?"

### Concrete result

In one validation case, three related firing conditions were delivered as a single grouped notification:

```text
[FIRING:3]
```

That is exactly the behavior I wanted. Three pieces of evidence, one operator interruption.

The alerts were not hidden. The notification preserved the multiplicity while collapsing redundant delivery.

---

## Change 2: use `group_wait` to absorb near-simultaneous fan-out

The scheduler-specific policy used:

```text
group_wait = 2m
```

A short group wait gives sibling alerts a chance to enter the same notification batch.

Too short and you still get fan-out. Too long and you delay the first useful signal.

Two minutes was a practical compromise for this workload because the alerts were scheduler-oriented rather than sub-second request-path failures.

The key is that `group_wait` is not a universal constant. It should reflect how quickly correlated alerts tend to arrive and how much initial notification latency the service can tolerate.

---

## Change 3: keep the group open long enough for related state changes

The policy used:

```text
group_interval = 30m
```

That meant changes to an existing alert group could be consolidated instead of immediately producing a new message for every small membership change.

The group interval is where you decide how often a *changed* group is worth interrupting a human again.

For scheduler failures, 30 minutes was acceptable because the first alert was already delivered after the 2-minute group wait. The 30-minute group interval affected subsequent grouped updates, not initial detection.

---

## Change 4: stop reminding humans about unchanged state every 30 minutes

This was the highest-leverage change.

The previous repeat behavior was:

```text
repeat_interval = 30m
```

That means a still-firing condition can keep reminding the team every half hour even when nothing has changed.

For the scoped scheduler route, I changed it to:

```text
repeat_interval = 720h
```

`720h` is 30 days.

This looks aggressive if you read it as "ignore the alert for 30 days." That is not what it means.

The first firing notification is still sent. Group membership changes can still produce updates according to `group_interval`. Resolved state is still part of the alert lifecycle. The long repeat interval only suppresses periodic reminders for an otherwise unchanged firing group.

In other words:

```text
new information      -> notify
changed group        -> notify on group policy
same state, no change -> do not nag every 30m
```

---

## Scope mattered more than the number

I did **not** apply the 30-day repeat interval to the global production policy.

The general production route retained a much shorter repeat cadence. The 30-day setting was applied only to the scheduler-specific route where the operational semantics supported it.

That distinction matters.

A common alert-fatigue failure mode is making a global policy less sensitive because one alert family is noisy. That trades local annoyance for global risk.

The safer pattern is:

```text
default policy
    |
    +-- ordinary production alerts -> normal repeat behavior
    |
    +-- scheduler alert family     -> scheduler-specific grouping/repeat behavior
```

Tune the route that is noisy, not the entire monitoring system.

---

## What I deliberately did not change

I did not use noise reduction as an excuse to weaken failure detection.

The change did not depend on:

- deleting scheduler failure rules
- raising error thresholds until alerts stopped
- increasing the evaluation interval to hours
- hiding alerts behind a long pending period
- disabling resolved-state handling
- globally suppressing production reminders

That separation is important because "fewer notifications" is not a reliability objective by itself.

The objective is **fewer low-information interruptions without losing actionable state transitions**.

---

## Validation

I used two levels of validation.

### 1. Alert-state validation

A controlled test rule was driven through:

```text
Normal -> Alerting -> Normal
```

This confirmed that changing notification behavior did not break state evaluation.

### 2. Workload validation

In one production scheduler validation window, the execution counters were:

```text
expected executions: 24
started:             24
successful:          24
errors:               0
```

The point of this check was not to prove the alerting policy from workload success. It was to make sure the monitoring change was not accidentally altering or misrepresenting the scheduler execution path.

---

## The Terraform/Grafana shape

The exact surrounding resource structure depends on your Grafana provider version, but the policy values looked like this conceptually:

```hcl
policy {
  group_by = [
    "grafana_folder",
    "environment",
    "alert_type",
    "application",
  ]

  group_wait      = "2m"
  group_interval  = "30m"
  repeat_interval = "720h"
}
```

The configuration itself is simple. The hard part is choosing labels that represent an operational incident rather than an implementation detail.

If `application` or `alert_type` is inconsistent across rules, grouping will fragment. Before changing policy, normalize the labels that define incident identity.

---

## Why this worked

The 35% reduction came from attacking notification multiplication rather than alert existence.

Think of notification volume as roughly:

```text
notifications
  = distinct alert groups
  x state changes worth delivering
  x repeat reminders
```

There are three levers:

1. reduce unnecessary group cardinality
2. batch correlated state changes
3. eliminate unchanged-state reminders

This tuning cycle used all three, while leaving the actual detection conditions intact.

---

## What can go wrong

### Over-grouping

If the grouping key is too broad, unrelated failures collapse into one notification and operators lose service identity.

Bad example:

```text
[environment]
```

One production notification could then contain failures from unrelated applications.

### Under-grouping

If the key contains high-cardinality or rule-specific labels, you are back to one notification per rule.

Typical offenders include instance IDs, pod names, request IDs, or overly specific alert names.

### Long repeat intervals on the wrong route

A 30-day repeat interval is reasonable only when the first notification and meaningful state changes are sufficient for that alert family.

Do not copy `720h` into a global policy.

### Measuring rule count instead of interruption count

Ninety-one rules can be perfectly manageable if they collapse into a small number of actionable notifications. Ten rules can be unbearable if each repeats every 30 minutes.

Rule count is an inventory metric. Notification count is an attention metric.

---

## A practical review sequence

When I review a noisy Grafana policy now, I work in this order:

1. **Identify the human interruption path.** Which contact point actually wakes or interrupts someone?
2. **Separate detection from delivery.** Confirm whether the rule is wrong or the notification semantics are wrong.
3. **Inspect grouping cardinality.** Ask what labels define one operational incident.
4. **Inspect repeat behavior.** Determine whether unchanged state is generating reminders without new information.
5. **Scope the fix.** Route-specific tuning is safer than weakening the default policy.
6. **Drive a controlled state transition.** Test `Normal -> Alerting -> Normal`.
7. **Compare notification volume.** Count delivered notifications over equivalent windows, not raw rule evaluations.
8. **Review missed-signal risk.** Verify that no meaningful transition was suppressed.

---

## The bigger SRE lesson

Alert fatigue is often treated as a threshold problem.

Sometimes it is. But a large class of alert fatigue is really a **cardinality and notification-state problem**.

Before weakening detection, inspect the path between "rule is firing" and "human gets interrupted."

In this case, changing four notification-policy dimensions was enough to cut alert noise by **35%**:

```text
group_by         -> operational identity
group_wait       -> 2m
group_interval   -> 30m
repeat_interval  -> 720h on the scoped scheduler route
```

The alerts still fired. The operators just stopped being told the same thing over and over again.

That is the kind of alert-noise reduction I want: less interruption, same signal.
