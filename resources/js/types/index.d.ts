import { LucideIcon } from 'lucide-react';
import type { Config } from 'ziggy-js';

export interface Auth {
    user: User;
}

export interface BreadcrumbItem {
    title: string;
    href: string;
}

export interface NavGroup {
    title: string;
    items: NavItem[];
}

export interface NavItem {
    title: string;
    href: string;
    icon?: LucideIcon | null;
    isActive?: boolean;
    items?: NavItem[];
}

export interface LayoutOption {
    name: string;
    description: string;
    component: string;
}

export interface LayoutConfig {
    current: string;
    options: Record<string, LayoutOption>;
}

export interface SharedData {
    name: string;
    quote: { message: string; author: string };
    auth: Auth;
    ziggy: Config & { location: string };
    sidebarOpen: boolean;
    layout: LayoutConfig;
    [key: string]: unknown;
}

export interface User {
    id: number;
    name: string;
    username: string;
    avatar?: string;
    created_at: string;
    updated_at: string;
    [key: string]: unknown; // This allows for additional properties...
}
