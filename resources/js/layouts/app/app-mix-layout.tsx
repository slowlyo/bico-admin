import { AppContent } from '@/components/app-content';
import { AppHeaderMix } from '@/components/app-header-mix';
import { AppShell } from '@/components/app-shell';
import { AppSidebar } from '@/components/app-sidebar';
import { AppSidebarHeader } from '@/components/app-sidebar-header';
import { type BreadcrumbItem } from '@/types';
import { type PropsWithChildren } from 'react';

/**
 * 混合布局组件
 *
 * 设计理念：
 * - 左侧：侧边栏提供主导航功能
 * - 顶部：显示用户信息和面包屑导航
 * - 桌面端和移动端都保持一致的布局结构
 * - 简洁高效的导航体验
 */
export default function AppMixLayout({ children, breadcrumbs = [] }: PropsWithChildren<{ breadcrumbs?: BreadcrumbItem[] }>) {
    return (
        <AppShell variant="sidebar">
            {/* 左侧：侧边栏导航 */}
            <AppSidebar />

            <AppContent variant="sidebar" className="overflow-x-hidden">
                {/* 顶部：用户信息和面包屑（桌面端） */}
                <div className="hidden lg:block">
                    <AppHeaderMix breadcrumbs={breadcrumbs} />
                </div>

                {/* 顶部：侧边栏头部（移动端，包含侧边栏切换和面包屑） */}
                <div className="lg:hidden">
                    <AppSidebarHeader breadcrumbs={breadcrumbs} />
                </div>

                {/* 主要内容区域 */}
                <div className="flex-1 px-6 py-6">
                    {children}
                </div>
            </AppContent>
        </AppShell>
    );
}
