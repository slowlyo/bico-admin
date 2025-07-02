import { type BreadcrumbItem, type SharedData, type User } from '@/types';
import { Head, Link, usePage } from '@inertiajs/react';
import { ArrowLeft, Calendar, Edit, Hash, User as UserIcon } from 'lucide-react';
import AppLayout from '@/layouts/app-layout';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';

interface Props {
    user: User;
}

export default function ShowUser() {
    const { user } = usePage<SharedData & Props>().props;

    const breadcrumbs: BreadcrumbItem[] = [
        { title: '仪表盘', href: route('dashboard') },
        { title: '系统管理', href: '#' },
        { title: '用户管理', href: route('admin.users.index') },
        { title: '用户详情', href: route('admin.users.show', user.id) },
    ];

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleString('zh-CN', {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
        });
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title={`用户详情 - ${user.name}`} />

            <div className="space-y-6">
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <Link href={route('admin.users.index')}>
                            <Button variant="outline" size="sm">
                                <ArrowLeft className="mr-2 h-4 w-4" />
                                返回列表
                            </Button>
                        </Link>
                        <div>
                            <h1 className="text-3xl font-bold tracking-tight">用户详情</h1>
                            <p className="text-muted-foreground">查看用户 {user.name} 的详细信息</p>
                        </div>
                    </div>
                    <Link href={route('admin.users.edit', user.id)}>
                        <Button>
                            <Edit className="mr-2 h-4 w-4" />
                            编辑用户
                        </Button>
                    </Link>
                </div>

                <div className="grid gap-6 md:grid-cols-2">
                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <UserIcon className="h-5 w-5" />
                                基本信息
                            </CardTitle>
                            <CardDescription>用户的基本资料信息</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <Hash className="h-4 w-4 text-muted-foreground" />
                                    <span className="text-sm font-medium">用户ID</span>
                                </div>
                                <Badge variant="outline">{user.id}</Badge>
                            </div>

                            <Separator />

                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <UserIcon className="h-4 w-4 text-muted-foreground" />
                                    <span className="text-sm font-medium">姓名</span>
                                </div>
                                <span className="text-sm">{user.name}</span>
                            </div>

                            <Separator />

                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-2">
                                    <UserIcon className="h-4 w-4 text-muted-foreground" />
                                    <span className="text-sm font-medium">用户名</span>
                                </div>
                                <Badge variant="secondary">@{user.username}</Badge>
                            </div>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle className="flex items-center gap-2">
                                <Calendar className="h-5 w-5" />
                                时间信息
                            </CardTitle>
                            <CardDescription>用户的创建和更新时间</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="space-y-2">
                                <div className="flex items-center gap-2">
                                    <Calendar className="h-4 w-4 text-muted-foreground" />
                                    <span className="text-sm font-medium">创建时间</span>
                                </div>
                                <p className="text-sm text-muted-foreground pl-6">
                                    {formatDate(user.created_at)}
                                </p>
                            </div>

                            <Separator />

                            <div className="space-y-2">
                                <div className="flex items-center gap-2">
                                    <Calendar className="h-4 w-4 text-muted-foreground" />
                                    <span className="text-sm font-medium">最后更新</span>
                                </div>
                                <p className="text-sm text-muted-foreground pl-6">
                                    {formatDate(user.updated_at)}
                                </p>
                            </div>
                        </CardContent>
                    </Card>
                </div>

                <Card>
                    <CardHeader>
                        <CardTitle>操作记录</CardTitle>
                        <CardDescription>用户的相关操作历史</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="text-center py-8 text-muted-foreground">
                            <Calendar className="mx-auto h-12 w-12 mb-4 opacity-50" />
                            <p>暂无操作记录</p>
                            <p className="text-sm">用户的操作历史将在此处显示</p>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </AppLayout>
    );
}
