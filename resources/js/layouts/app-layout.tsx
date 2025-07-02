import AppHeaderLayout from '@/layouts/app/app-header-layout';
import AppSidebarLayout from '@/layouts/app/app-sidebar-layout';
import { type BreadcrumbItem, type SharedData } from '@/types';
import { usePage } from '@inertiajs/react';
import { type ReactNode } from 'react';

interface AppLayoutProps {
    children: ReactNode;
    breadcrumbs?: BreadcrumbItem[];
}

export default function AppLayout({ children, breadcrumbs, ...props }: AppLayoutProps) {
    const { layout } = usePage<SharedData>().props;

    // 根据配置选择布局组件
    const LayoutComponent = getLayoutComponent(layout.current);

    return (
        <LayoutComponent breadcrumbs={breadcrumbs} {...props}>
            {children}
        </LayoutComponent>
    );
}

/**
 * 根据布局类型获取对应的布局组件
 */
function getLayoutComponent(layoutType: string) {
    switch (layoutType) {
        case 'header':
            return AppHeaderLayout;
        case 'sidebar':
        default:
            return AppSidebarLayout;
    }
}
