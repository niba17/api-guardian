import React from "react";
import {
  LayoutDashboard,
  ShieldAlert,
  Activity,
  Settings,
  Menu,
  Shield,
  LogOut, // Pastikan LogOut terimport
} from "lucide-react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";

interface MainLayoutProps {
  children: React.ReactNode;
}

export default function MainLayout({ children }: MainLayoutProps) {
  const location = useLocation();
  const { logout } = useAuth();
  const navigate = useNavigate(); // 👈 PINDAHKAN KE SINI (Top Level)

  const handleLogout = () => {
    logout(); // 👈 Panggil fungsi logout dari context
    navigate("/login");
  };

  const menuItems = [
    { icon: LayoutDashboard, label: "Overview", path: "/" },
    { icon: ShieldAlert, label: "Threats & WAF", path: "/threats" },
    { icon: Activity, label: "Traffic Monitor", path: "/traffic" },
    { icon: Settings, label: "Configuration", path: "/settings" },
  ];

  return (
    <div className="min-h-screen bg-guardian-dark text-slate-300 font-sans selection:bg-guardian-primary selection:text-white">
      {/* SIDEBAR */}
      <aside className="fixed left-0 top-0 h-screen w-64 bg-guardian-card border-r border-slate-800 hidden md:flex flex-col z-50">
        {/* Logo Area */}
        <div className="p-6 border-b border-slate-800 flex items-center gap-3">
          <div className="bg-guardian-primary/20 p-2 rounded-lg">
            <Shield className="text-guardian-primary w-6 h-6" />
          </div>
          <div>
            <h1 className="font-bold text-white tracking-wider">GUARDIAN</h1>
            <p className="text-[10px] text-slate-500 font-mono">
              API SECURITY GATEWAY
            </p>
          </div>
        </div>

        {/* Navigation */}
        <nav className="flex-1 p-4 space-y-2 overflow-y-auto">
          {menuItems.map((item) => {
            const isActive = location.pathname === item.path;

            return (
              <Link
                key={item.label}
                to={item.path}
                className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 group ${
                  isActive
                    ? "bg-guardian-primary text-white shadow-lg shadow-guardian-primary/20"
                    : "hover:bg-slate-800 hover:text-white"
                }`}
              >
                <item.icon
                  size={20}
                  className={
                    isActive
                      ? "text-white"
                      : "text-slate-500 group-hover:text-white"
                  }
                />
                <span className="font-medium text-sm">{item.label}</span>
                {isActive && (
                  <div className="ml-auto w-1.5 h-1.5 rounded-full bg-white animate-pulse" />
                )}
              </Link>
            );
          })}

          {/* Tombol Logout */}
          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 hover:bg-red-500/10 hover:text-red-400 text-slate-500 group mt-4 border-t border-slate-800"
          >
            <LogOut size={20} className="group-hover:text-red-400" />
            <span className="font-medium text-sm">Logout</span>
          </button>
        </nav>

        {/* Footer Sidebar */}
        <div className="p-4 border-t border-slate-800">
          <div className="bg-slate-900/50 rounded-xl p-4 border border-slate-800">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
              <span className="text-xs font-bold text-green-500">
                SYSTEM STABLE
              </span>
            </div>
            <p className="text-[10px] text-slate-500">v1.0.2 • Build 2026</p>
          </div>
        </div>
      </aside>

      {/* MAIN CONTENT */}
      <main className="md:ml-64 min-h-screen relative">
        {/* Top Mobile Bar */}
        <div className="md:hidden h-16 bg-guardian-card border-b border-slate-800 flex items-center justify-between px-4 sticky top-0 z-40">
          <span className="font-bold text-white">GUARDIAN</span>
          <button className="p-2">
            <Menu />
          </button>
        </div>

        <div className="max-w-7xl mx-auto p-6 lg:p-8">{children}</div>
      </main>
    </div>
  );
}
