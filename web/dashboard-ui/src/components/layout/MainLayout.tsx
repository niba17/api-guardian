import React from "react";
import {
  LayoutDashboard,
  ShieldAlert,
  Activity,
  Settings,
  Menu,
  Shield,
  LogOut,
  Clock, // 👈 Tambahkan icon Clock untuk estetika
} from "lucide-react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import TimezoneSelector from "../TimezoneSelector"; // 👈 Tanpa kurung kurawal {}

interface MainLayoutProps {
  children: React.ReactNode;
}

export default function MainLayout({ children }: MainLayoutProps) {
  const location = useLocation();
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
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
      {/* SIDEBAR (Desktop) */}
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

          <button
            onClick={handleLogout}
            className="w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 hover:bg-red-500/10 hover:text-red-400 text-slate-500 group mt-4 border-t border-slate-800 pt-6"
          >
            <LogOut size={20} />
            <span className="font-medium text-sm">Logout</span>
          </button>
        </nav>

        {/* FOOTER SIDEBAR (Pusat Kendali Waktu) */}
        <div className="p-4 border-t border-slate-800 space-y-4">
          {/* TAMPILAN TIMEZONE SELECTOR */}
          <div className="px-2">
            <div className="flex items-center gap-2 mb-2 text-slate-500">
              <Clock size={12} />
              <span className="text-[10px] uppercase font-bold tracking-widest">
                Display Region
              </span>
            </div>
            <TimezoneSelector />
          </div>

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
        {/* Top Mobile Bar (Dengan Selector untuk Mobile) */}
        <div className="md:hidden h-16 bg-guardian-card border-b border-slate-800 flex items-center justify-between px-4 sticky top-0 z-40">
          <div className="flex items-center gap-2">
            <Shield className="text-guardian-primary w-5 h-5" />
            <span className="font-bold text-white text-sm">GUARDIAN</span>
          </div>

          <div className="flex items-center gap-3">
            <TimezoneSelector /> {/* 👈 Muncul juga di mobile bar */}
            <button className="p-2 text-slate-400">
              <Menu size={20} />
            </button>
          </div>
        </div>

        <div className="max-w-7xl mx-auto p-6 lg:p-8">{children}</div>
      </main>
    </div>
  );
}
