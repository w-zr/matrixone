# MERGE Scheduling

## Scheduler

The scheduler traverses all objects. These objects are sequentially checked by each policy. If an object satisfies a
policy, it will be skipped by subsequent policies. Every merge operation results in objects that belong to the same
segment.

## Policies

### Basic Policy

The basic policy is used to merge smaller objects into larger ones. If an object is smaller than 110MB, it will be
included in this policy and merged into a size of no less than 128MB during subsequent execution. However, a table can
only merge 16 objects at a time or objects with a total size less than 256MB.

### Compact Policy

The compaction policy applies tombstones to objects. This policy records tombstones that are larger than 110MB or have
existed for more than 10 minutes and are larger than 10MB. Then, the policy checks which objects contain data within
these tombstones and merges these objects to reduce the size of both tombstones and objects.

### Overlap Policy

The Overlap policy first groups all objects by the number of objects in each segment. Then, within each group, it
checks whether the zonemaps of the objects have overlaps. If there is an overlap, then these objects are merged into
a new segment.

### Tombstone Policy

The tombstone policy aims to merge multiple tombstones into a single tombstone.

## Memory check

Merging a large number of objects can consume a significant amount of memory. If there is no limit on the amount of data
merged at one time, TN is very likely to run out of memory. For estimating memory consumption, we base our judgment on
the total number of rows merged. The current estimate is that each row occupies 30 bytes of memory. The current limit on
memory consumption for merges is controlled to be 2/3 of the available memory. If a merge operation requires more
memory than 2/3 of the current available memory, then that merge will not be performed.
