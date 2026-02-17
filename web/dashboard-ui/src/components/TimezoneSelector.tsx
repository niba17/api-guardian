// web\dashboard-ui\src\components\TimezoneSelector.tsx
import { getCurrentTimezone, setGlobalTimezone } from "../utils/time";

// WAJIB: Pastikan ada kata 'default' di sini!
export default function TimezoneSelector() {
  const tzs = [
    { label: "WIB (Jakarta)", value: "Asia/Jakarta" },
    { label: "WITA (Kupang/Mks)", value: "Asia/Makassar" },
    { label: "WIT (Jayapura)", value: "Asia/Jayapura" },
    { label: "UTC (Global)", value: "UTC" },
  ];

  return (
    <select
      className="bg-slate-800 text-slate-200 text-[10px] p-1 rounded border border-slate-700 outline-none hover:border-guardian-primary transition-colors cursor-pointer"
      value={getCurrentTimezone()}
      onChange={(e) => setGlobalTimezone(e.target.value)}
    >
      {tzs.map((tz) => (
        <option key={tz.value} value={tz.value}>
          {tz.label}
        </option>
      ))}
    </select>
  );
}
