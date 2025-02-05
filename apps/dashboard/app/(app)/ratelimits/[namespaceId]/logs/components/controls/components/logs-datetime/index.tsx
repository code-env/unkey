import { cn } from "@/lib/utils";
import { Calendar } from "@unkey/icons";
import { Button } from "@unkey/ui";
import { useState } from "react";
import { useFilters } from "../../../../hooks/use-filters";
import { DatetimePopover } from "./components/datetime-popover";

export const LogsDateTime = () => {
  const [title, setTitle] = useState<string>("Last 12 hours");
  const { filters, updateFilters } = useFilters();
  const timeValues = filters
    .filter((f) => ["startTime", "endTime", "since"].includes(f.field))
    .reduce(
      (acc, f) => ({
        // biome-ignore lint/performance/noAccumulatingSpread: safe to spread
        ...acc,
        [f.field]: f.value,
      }),
      {},
    );

  return (
    <DatetimePopover
      initialTimeValues={timeValues}
      onDateTimeChange={(startTime, endTime, since) => {
        const activeFilters = filters.filter(
          (f) => !["endTime", "startTime", "since"].includes(f.field),
        );
        if (since !== undefined) {
          updateFilters([
            ...activeFilters,
            {
              field: "since",
              value: since,
              id: crypto.randomUUID(),
              operator: "is",
            },
          ]);
          return;
        }
        if (since === undefined && startTime) {
          activeFilters.push({
            field: "startTime",
            value: startTime,
            id: crypto.randomUUID(),
            operator: "is",
          });
          if (endTime) {
            activeFilters.push({
              field: "endTime",
              value: endTime,
              id: crypto.randomUUID(),
              operator: "is",
            });
          }
        }
        updateFilters(activeFilters);
      }}
      initialTitle={title}
      onSuggestionChange={(newTitle) => {
        setTitle(newTitle);
      }}
    >
      <div className="group">
        <Button
          variant="ghost"
          className={cn(
            "group-data-[state=open]:bg-gray-4 px-2",
            title !== "Last 12 hours" ? "bg-gray-4" : "",
          )}
          aria-label="Filter logs by time"
          aria-haspopup="true"
          title="Press 'T' to toggle filters"
        >
          <Calendar className="text-gray-9 size-4" />
          <span className="text-gray-12 font-medium text-[13px]">{title}</span>
        </Button>
      </div>
    </DatetimePopover>
  );
};
