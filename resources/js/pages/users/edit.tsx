import { type BreadcrumbItem, type SharedData, type User } from '@/types';
import { Head, Link, useForm, usePage } from '@inertiajs/react';
import { ArrowLeft, Save, User as UserIcon } from 'lucide-react';
import { FormEventHandler, useState } from 'react';
import AppLayout from '@/layouts/app-layout';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Checkbox } from '@/components/ui/checkbox';
import { Input } from '@/components/ui/input';
import InputError from '@/components/input-error';
import { Label } from '@/components/ui/label';

interface Props {
    user: User;
}

type EditUserForm = {
    name: string;
    password: string;
    password_confirmation: string;
};

export default function EditUser() {
    const { user } = usePage<SharedData & Props>().props;
    const [resetPassword, setResetPassword] = useState(false);

    const { data, setData, patch, errors, processing } = useForm<EditUserForm>({
        name: user.name,
        password: '',
        password_confirmation: '',
    });

    const breadcrumbs: BreadcrumbItem[] = [
        { title: '仪表盘', href: route('dashboard') },
        { title: '系统管理', href: '#' },
        { title: '用户管理', href: route('admin.users.index') },
        { title: '编辑用户', href: route('admin.users.edit', user.id) },
    ];

    const submit: FormEventHandler = (e) => {
        e.preventDefault();

        if (resetPassword) {
            setData({
                name: data.name,
                password: data.password,
                password_confirmation: data.password_confirmation,
            });
        } else {
            setData({
                name: data.name,
                password: '',
                password_confirmation: '',
            });
        }

        patch(route('admin.users.update', user.id));
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title={`编辑用户 - ${user.name}`} />

            <div className="space-y-6">
                <div className="flex items-center gap-4">
                    <Link href={route('admin.users.index')}>
                        <Button variant="outline" size="sm">
                            <ArrowLeft className="mr-2 h-4 w-4" />
                            返回列表
                        </Button>
                    </Link>
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">编辑用户</h1>
                        <p className="text-muted-foreground">修改用户 {user.name} 的信息</p>
                    </div>
                </div>

                <Card className="max-w-2xl">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <UserIcon className="h-5 w-5" />
                            用户信息
                        </CardTitle>
                        <CardDescription>修改用户的基本信息</CardDescription>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={submit} className="space-y-6">
                            <div className="grid gap-2">
                                <Label htmlFor="username">用户名</Label>
                                <Input
                                    id="username"
                                    type="text"
                                    value={user.username}
                                    disabled
                                    className="bg-muted"
                                />
                                <p className="text-sm text-muted-foreground">
                                    用户名不可更改
                                </p>
                            </div>

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

                            <div className="space-y-4">
                                <div className="flex items-center space-x-2">
                                    <Checkbox
                                        id="reset-password"
                                        checked={resetPassword}
                                        onCheckedChange={(checked) => {
                                            setResetPassword(checked as boolean);
                                            if (!checked) {
                                                setData('password', '');
                                                setData('password_confirmation', '');
                                            }
                                        }}
                                    />
                                    <Label htmlFor="reset-password">重置密码</Label>
                                </div>

                                {resetPassword && (
                                    <div className="space-y-4 pl-6 border-l-2 border-muted">
                                        <div className="grid gap-2">
                                            <Label htmlFor="password">新密码 *</Label>
                                            <Input
                                                id="password"
                                                type="password"
                                                value={data.password}
                                                onChange={(e) => setData('password', e.target.value)}
                                                required={resetPassword}
                                                autoComplete="new-password"
                                                placeholder="请输入新密码（至少6个字符）"
                                            />
                                            <InputError message={errors.password} />
                                        </div>

                                        <div className="grid gap-2">
                                            <Label htmlFor="password_confirmation">确认新密码 *</Label>
                                            <Input
                                                id="password_confirmation"
                                                type="password"
                                                value={data.password_confirmation}
                                                onChange={(e) => setData('password_confirmation', e.target.value)}
                                                required={resetPassword}
                                                autoComplete="new-password"
                                                placeholder="请再次输入新密码"
                                            />
                                            <InputError message={errors.password_confirmation} />
                                        </div>
                                    </div>
                                )}
                            </div>

                            <div className="flex items-center gap-4 pt-4">
                                <Button type="submit" disabled={processing}>
                                    <Save className="mr-2 h-4 w-4" />
                                    {processing ? '保存中...' : '保存更改'}
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
