import { requestClient } from '#/api/request';

export interface DashboardMetric {
  key: string;
  label: string;
  status: string;
  unit: string;
  value: number;
}

export interface DashboardOverview {
  database: {
    driver: string;
    idle: number;
    inUse: number;
    maxOpenConnections: number;
    openConnections: number;
    waitCount: number;
  };
  monitor: {
    collectedAt: string;
    metrics: DashboardMetric[];
  };
  runtime: {
    allocMb: number;
    cpuCores: number;
    gcCycles: number;
    goroutines: number;
    sysMb: number;
  };
  server: {
    arch: string;
    goVersion: string;
    hostname: string;
    mode: string;
    os: string;
    port: number;
    startedAt: string;
    uptimeSeconds: number;
  };
}

// 工作台首屏数据。接口只汇总运行状态，不返回业务计数。
export function getDashboardOverviewApi() {
  return requestClient.get<DashboardOverview>('/dashboard/overview');
}
