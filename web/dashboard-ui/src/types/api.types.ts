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
}

export interface SystemStats {
  total_requests: number;
  blocked_requests: number;
  unique_ips: number; // 👈 INI YANG TADI HILANG/SALAH
  avg_latency: string;
}
