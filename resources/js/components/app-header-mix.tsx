import { Breadcrumbs } from '@/components/breadcrumbs';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { UserMenuContent } from '@/components/user-menu-content';
import { useInitials } from '@/hooks/use-initials';
import { type BreadcrumbItem, type SharedData } from '@/types';
import { usePage } from '@inertiajs/react';



interface AppHeaderMixProps {
    breadcrumbs?: BreadcrumbItem[];
}

/**
 * 混合布局专用的顶部栏组件
 *
 * 设计理念：
 * - 左侧：侧边栏提供主导航
 * - 顶部：只显示用户信息和面包屑导航
 * - 简洁的顶部栏，不包含主导航菜单
 */
export function AppHeaderMix({ breadcrumbs = [] }: AppHeaderMixProps) {
    const { auth } = usePage<SharedData>().props;
    const getInitials = useInitials();

    return (
        <>
            {/* 顶部栏：用户信息区域 */}
            <div className="border-b border-sidebar-border/80">
                <div className="flex h-16 items-center justify-end px-6">
                    {/* 用户菜单 */}
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" className="size-10 rounded-full p-1">
                                <Avatar className="size-8 overflow-hidden rounded-full">
                                    <AvatarImage src={auth.user.avatar} alt={auth.user.name} />
                                    <AvatarFallback className="rounded-lg bg-neutral-200 text-black dark:bg-neutral-700 dark:text-white">
                                        {getInitials(auth.user.name)}
                                    </AvatarFallback>
                                </Avatar>
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent className="w-56" align="end">
                            <UserMenuContent user={auth.user} />
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </div>

            {/* 面包屑导航区域 */}
            {breadcrumbs.length > 1 && (
                <div className="flex w-full border-b border-sidebar-border/70">
                    <div className="w-full px-6 py-3">
                        <Breadcrumbs breadcrumbs={breadcrumbs} />
                    </div>
                </div>
            )}
        </>
    );
}
