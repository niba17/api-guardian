import {
  LayoutDashboard,
  ShieldAlert,
  Activity,
  FileText,
  Settings,
  ShieldCheck,
} from "lucide-react";

const menuItems = [
  { icon: LayoutDashboard, label: "Overview", active: true },
  { icon: ShieldAlert, label: "Threats & WAF", active: false },
  { icon: Activity, label: "Traffic Monitor", active: false },
  { icon: FileText, label: "Audit Logs", active: false },
  { icon: Settings, label: "Configuration", active: false },
];

export default function Sidebar() {
  return (
    <aside className="fixed left-0 top-0 h-screen w-64 bg-guardian-card border-r border-slate-800 flex flex-col">
      <div className="p-6 flex items-center gap-3 border-b border-slate-800">
        <div className="w-10 h-10 bg-guardian-primary/20 rounded-xl flex items-center justify-center border border-guardian-primary/30">
          <ShieldCheck className="text-guardian-primary" size={24} />
        </div>
        <span className="font-bold text-xl tracking-tight">GUARDIAN</span>
      </div>

      <nav className="flex-1 p-4 space-y-2 mt-4">
        {menuItems.map((item, idx) => (
          <button
            key={idx}
            className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 group ${
              item.active
                ? "bg-guardian-primary/10 text-guardian-primary border border-guardian-primary/20"
                : "hover:bg-slate-800/50 text-slate-400 hover:text-white"
            }`}
          >
            <item.icon size={20} />
            <span className="font-medium">{item.label}</span>
          </button>
        ))}
      </nav>

      <div className="p-4 border-t border-slate-800">
        <div className="bg-slate-900/50 p-4 rounded-2xl border border-slate-800">
          <p className="text-xs text-slate-500 font-mono mb-2 uppercase tracking-widest">
            System Status
          </p>
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-guardian-primary animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.6)]"></span>
            <span className="text-sm font-semibold text-guardian-primary">
              LIVE MONITORING
            </span>
          </div>
        </div>
      </div>
    </aside>
  );
}
