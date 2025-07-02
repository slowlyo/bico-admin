import { type BreadcrumbItem } from '@/types';
import { Head, Link, useForm } from '@inertiajs/react';
import { ArrowLeft, Save, User } from 'lucide-react';
import { FormEventHandler } from 'react';
import AppLayout from '@/layouts/app-layout';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import InputError from '@/components/input-error';
import { Label } from '@/components/ui/label';

type CreateUserForm = {
    name: string;
    username: string;
    password: string;
    password_confirmation: string;
};

export default function CreateUser() {
    const { data, setData, post, errors, processing } = useForm<CreateUserForm>({
        name: '',
        username: '',
        password: '',
        password_confirmation: '',
    });

    const breadcrumbs: BreadcrumbItem[] = [
        { title: '仪表盘', href: route('dashboard') },
        { title: '系统管理', href: '#' },
        { title: '用户管理', href: route('admin.users.index') },
        { title: '创建用户', href: route('admin.users.create') },
    ];

    const submit: FormEventHandler = (e) => {
        e.preventDefault();
        post(route('admin.users.store'));
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="创建用户" />

            <div className="space-y-6">
                <div className="flex items-center gap-4">
                    <Link href={route('admin.users.index')}>
                        <Button variant="outline" size="sm">
                            <ArrowLeft className="mr-2 h-4 w-4" />
                            返回列表
                        </Button>
                    </Link>
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">创建用户</h1>
                        <p className="text-muted-foreground">添加新的系统用户</p>
                    </div>
                </div>

                <Card className="max-w-2xl">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <User className="h-5 w-5" />
                            用户信息
                        </CardTitle>
                        <CardDescription>请填写用户的基本信息</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={submit} className="space-y-6">
                            <div className="grid gap-2">
                                <Label htmlFor="name">姓名 *</Label>
                                <Input
                                    id="name"
                                    type="text"
                                    value={data.name}
                                    onChange={(e) => setData('name', e.target.value)}
                                    required
                                    autoComplete="name"
                                    placeholder="请输入用户姓名"
                                />
                                <InputError message={errors.name} />
                            </div>

                            <div className="grid gap-2">
                                <Label htmlFor="username">用户名 *</Label>
                                <Input
                                    id="username"
                                    type="text"
                                    value={data.username}
                                    onChange={(e) => setData('username', e.target.value)}
                                    required
                                    autoComplete="username"
                                    placeholder="请输入用户名（只能包含字母、数字、破折号和下划线）"
                                />
                                <InputError message={errors.username} />
                                <p className="text-sm text-muted-foreground">
                                    用户名将用于登录，创建后不可更改
                                </p>
                            </div>

                            <div className="grid gap-2">
                                <Label htmlFor="password">密码 *</Label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={data.password}
                                    onChange={(e) => setData('password', e.target.value)}
                                    required
                                    autoComplete="new-password"
                                    placeholder="请输入密码（至少6个字符）"
                                />
                                <InputError message={errors.password} />
                            </div>

                            <div className="grid gap-2">
                                <Label htmlFor="password_confirmation">确认密码 *</Label>
                                <Input
                                    id="password_confirmation"
                                    type="password"
                                    value={data.password_confirmation}
                                    onChange={(e) => setData('password_confirmation', e.target.value)}
                                    required
                                    autoComplete="new-password"
                                    placeholder="请再次输入密码"
                                />
                                <InputError message={errors.password_confirmation} />
                            </div>

                            <div className="flex items-center gap-4 pt-4">
                                <Button type="submit" disabled={processing}>
                                    <Save className="mr-2 h-4 w-4" />
                                    {processing ? '创建中...' : '创建用户'}
                                </Button>
                                <Link href={route('admin.users.index')}>
                                    <Button type="button" variant="outline">
                                        取消
                                    </Button>
                                </Link>
                            </div>
                        </form>
                    </CardContent>
                </Card>
            </div>
        </AppLayout>
    );
}
