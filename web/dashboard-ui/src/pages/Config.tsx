import { useState, useEffect } from "react";
import MainLayout from "../components/layout/MainLayout";
import Switch from "../components/ui/Switch";
import SettingsSkeleton from "../features/config/components/ConfigSkeleton";
import { Shield, Globe, Save, Activity } from "lucide-react";

export default function Config() {
  const [loading, setLoading] = useState(true);
  const [config, setConfig] = useState({
    wafEnabled: true,
    blockTor: false,
    rateLimit: true,
    geoBlock: false,
    strictMode: false,
  });

  useEffect(() => {
    const timer = setTimeout(() => setLoading(false), 800);
    return () => clearTimeout(timer);
  }, []);

  const handleToggle = (key: keyof typeof config) => {
    setConfig((prev) => ({ ...prev, [key]: !prev[key] }));
  };

  return (
    <MainLayout>
      {/* HEADER: Kloning Traffic.tsx */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white mb-2 flex items-center gap-2">
          {/* <Server className="text-guardian-primary" /> */}
          Configuration
          {loading && (
            <Activity
              className="animate-spin text-guardian-warning"
              size={20}
            />
          )}
        </h1>
        <p className="text-slate-400">
          Manage security protocols and firewall rules.
        </p>
      </div>

      {loading ? (
        <SettingsSkeleton />
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div className="bg-guardian-card border border-slate-800 rounded-2xl p-6 h-fit">
            <div className="flex items-center gap-3 mb-6 border-b border-slate-800 pb-4">
              <Shield className="text-emerald-400" />
              <h2 className="text-xl font-bold text-white">Core Protection</h2>
            </div>
            <div className="space-y-4">
              <Switch
                label="Active WAF"
                checked={config.wafEnabled}
                onChange={() => handleToggle("wafEnabled")}
              />
              <Switch
                label="Rate Limiting"
                checked={config.rateLimit}
                onChange={() => handleToggle("rateLimit")}
              />
              <Switch
                label="Paranoid Mode"
                checked={config.strictMode}
                onChange={() => handleToggle("strictMode")}
              />
            </div>
          </div>

          <div className="bg-guardian-card border border-slate-800 rounded-2xl p-6 h-fit">
            <div className="flex items-center gap-3 mb-6 border-b border-slate-800 pb-4">
              <Globe className="text-blue-400" />
              <h2 className="text-xl font-bold text-white">Network Rules</h2>
            </div>
            <div className="space-y-4">
              <Switch
                label="Block Tor"
                checked={config.blockTor}
                onChange={() => handleToggle("blockTor")}
              />
              <Switch
                label="Geo-Blocking"
                checked={config.geoBlock}
                onChange={() => handleToggle("geoBlock")}
              />
              <button className="mt-6 flex items-center justify-center gap-2 w-full bg-guardian-primary hover:bg-emerald-600 text-white font-bold py-3 rounded-xl transition-all">
                <Save size={18} /> Save Changes
              </button>
            </div>
          </div>
        </div>
      )}
    </MainLayout>
  );
}
