export interface SecurityLog {
  id: string;
  timestamp: string;
  ip: string;
  method: string;
  path: string;
  status: number;
  latency: number;
  country: string;
  city: string;
  is_blocked: boolean;
  threat_type?: string;
  threat_details?: string;
  user_agent?: string;
  browser?: string;
  os?: string;
  is_bot?: boolean;
  body?: string;
}

export interface SystemStats {
  total_requests: number;
  total_blocked: number;
  unique_ips: number; // 👈 INI YANG TADI HILANG/SALAH
  avg_latency: string;
}
