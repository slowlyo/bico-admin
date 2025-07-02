import { type BreadcrumbItem, type SharedData, type User } from '@/types';
import { Head, Link, router, usePage } from '@inertiajs/react';
import { Edit, Eye, Plus, Search, Trash2, Users } from 'lucide-react';
import { FormEventHandler, useState } from 'react';
import AppLayout from '@/layouts/app-layout';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

interface PaginationLink {
    url: string | null;
    label: string;
    active: boolean;
}

interface PaginatedUsers {
    data: User[];
    current_page: number;
    last_page: number;
    per_page: number;
    total: number;
    links: PaginationLink[];
}

interface Props {
    users: PaginatedUsers;
    filters: {
        search?: string;
        sort?: string;
        direction?: string;
    };
}

export default function UsersIndex() {
    const { users, filters } = usePage<SharedData & Props>().props;
    const [search, setSearch] = useState(filters.search || '');
    const [deleteUser, setDeleteUser] = useState<User | null>(null);
    const [isDeleting, setIsDeleting] = useState(false);

    const breadcrumbs: BreadcrumbItem[] = [
        { title: '仪表盘', href: route('dashboard') },
        { title: '系统管理', href: '#' },
        { title: '用户管理', href: route('admin.users.index') },
    ];

    const handleSearch: FormEventHandler = (e) => {
        e.preventDefault();
        router.get(route('admin.users.index'), { search }, { preserveState: true });
    };

    const handleSort = (field: string) => {
        const direction = filters.sort === field && filters.direction === 'asc' ? 'desc' : 'asc';
        router.get(route('admin.users.index'), { ...filters, sort: field, direction }, { preserveState: true });
    };

    const handleDelete = async () => {
        if (!deleteUser) return;

        setIsDeleting(true);
        try {
            await router.delete(route('admin.users.destroy', deleteUser.id));
            setDeleteUser(null);
        } finally {
            setIsDeleting(false);
        }
    };

    const getSortIcon = (field: string) => {
        if (filters.sort !== field) return null;
        return filters.direction === 'asc' ? '↑' : '↓';
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="用户管理" />

            <div className="space-y-6">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">用户管理</h1>
                        <p className="text-muted-foreground">管理系统用户账户</p>
                    </div>
                    <Link href={route('admin.users.create')}>
                        <Button>
                            <Plus className="mr-2 h-4 w-4" />
                            创建用户
                        </Button>
                    </Link>
                </div>

                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Users className="h-5 w-5" />
                            用户列表
                        </CardTitle>
                        <CardDescription>共 {users.total} 个用户</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <div className="mb-4">
                            <form onSubmit={handleSearch} className="flex gap-2">
                                <div className="relative flex-1">
                                    <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                                    <Input
                                        placeholder="搜索用户名或姓名..."
                                        value={search}
                                        onChange={(e) => setSearch(e.target.value)}
                                        className="pl-10"
                                    />
                                </div>
                                <Button type="submit">搜索</Button>
                                {filters.search && (
                                    <Button
                                        type="button"
                                        variant="outline"
                                        onClick={() => {
                                            setSearch('');
                                            router.get(route('admin.users.index'));
                                        }}
                                    >
                                        清除
                                    </Button>
                                )}
                            </form>
                        </div>

                        <div className="rounded-md border">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>ID</TableHead>
                                        <TableHead
                                            className="cursor-pointer hover:bg-muted/50"
                                            onClick={() => handleSort('name')}
                                        >
                                            姓名 {getSortIcon('name')}
                                        </TableHead>
                                        <TableHead
                                            className="cursor-pointer hover:bg-muted/50"
                                            onClick={() => handleSort('username')}
                                        >
                                            用户名 {getSortIcon('username')}
                                        </TableHead>
                                        <TableHead
                                            className="cursor-pointer hover:bg-muted/50"
                                            onClick={() => handleSort('created_at')}
                                        >
                                            创建时间 {getSortIcon('created_at')}
                                        </TableHead>
                                        <TableHead className="text-right">操作</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {users.data.length === 0 ? (
                                        <TableRow>
                                            <TableCell colSpan={5} className="text-center py-8 text-muted-foreground">
                                                {filters.search ? '没有找到匹配的用户' : '暂无用户数据'}
                                            </TableCell>
                                        </TableRow>
                                    ) : (
                                        users.data.map((user) => (
                                            <TableRow key={user.id}>
                                                <TableCell className="font-medium">{user.id}</TableCell>
                                                <TableCell>{user.name}</TableCell>
                                                <TableCell>
                                                    <Badge variant="secondary">@{user.username}</Badge>
                                                </TableCell>
                                                <TableCell>{new Date(user.created_at).toLocaleDateString('zh-CN')}</TableCell>
                                                <TableCell className="text-right">
                                                    <div className="flex justify-end gap-2">
                                                        <Link href={route('admin.users.show', user.id)}>
                                                            <Button variant="ghost" size="sm">
                                                                <Eye className="h-4 w-4" />
                                                            </Button>
                                                        </Link>
                                                        <Link href={route('admin.users.edit', user.id)}>
                                                            <Button variant="ghost" size="sm">
                                                                <Edit className="h-4 w-4" />
                                                            </Button>
                                                        </Link>
                                                        <Button
                                                            variant="ghost"
                                                            size="sm"
                                                            onClick={() => setDeleteUser(user)}
                                                            className="text-destructive hover:text-destructive"
                                                        >
                                                            <Trash2 className="h-4 w-4" />
                                                        </Button>
                                                    </div>
                                                </TableCell>
                                            </TableRow>
                                        ))
                                    )}
                                </TableBody>
                            </Table>
                        </div>

                        {/* 分页 */}
                        {users.last_page > 1 && (
                            <div className="mt-4 flex items-center justify-between">
                                <div className="text-sm text-muted-foreground">
                                    显示第 {(users.current_page - 1) * users.per_page + 1} 到{' '}
                                    {Math.min(users.current_page * users.per_page, users.total)} 条，共 {users.total} 条记录
                                </div>
                                <div className="flex gap-2">
                                    {users.links.map((link, index) => (
                                        <Button
                                            key={index}
                                            variant={link.active ? 'default' : 'outline'}
                                            size="sm"
                                            disabled={!link.url}
                                            onClick={() => link.url && router.get(link.url)}
                                            dangerouslySetInnerHTML={{ __html: link.label }}
                                        />
                                    ))}
                                </div>
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>

            {/* 删除确认对话框 */}
            <Dialog open={!!deleteUser} onOpenChange={() => setDeleteUser(null)}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>确认删除用户</DialogTitle>
                        <DialogDescription>
                            您确定要删除用户 <strong>{deleteUser?.name}</strong> (@{deleteUser?.username}) 吗？
                            此操作无法撤销。
                        </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setDeleteUser(null)}>
                            取消
                        </Button>
                        <Button variant="destructive" onClick={handleDelete} disabled={isDeleting}>
                            {isDeleting ? '删除中...' : '确认删除'}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </AppLayout>
    );
}
