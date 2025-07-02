import { type BreadcrumbItem, type SharedData } from '@/types';
import { Transition } from '@headlessui/react';
import { Head, useForm, usePage } from '@inertiajs/react';
import { FormEventHandler } from 'react';


import HeadingSmall from '@/components/heading-small';
import InputError from '@/components/input-error';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import AppLayout from '@/layouts/app-layout';
import SettingsLayout from '@/layouts/settings/layout';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: '个人资料设置',
        href: '/settings/profile',
    },
];

type ProfileForm = {
    name: string;
};

export default function Profile() {
    const { auth } = usePage<SharedData>().props;

    const { data, setData, patch, errors, processing, recentlySuccessful } = useForm<Required<ProfileForm>>({
        name: auth.user.name,
    });

    const submit: FormEventHandler = (e) => {
        e.preventDefault();

        patch(route('profile.update'), {
            preserveScroll: true,
        });
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="个人资料设置" />

            <SettingsLayout>
                <div className="space-y-6">
                    <HeadingSmall title="个人资料信息" description="更新您的姓名信息" />

                    <form onSubmit={submit} className="space-y-6">
                        <div className="grid gap-2">
                            <Label htmlFor="name">姓名</Label>

                            <Input
                                id="name"
                                className="mt-1 block w-full"
                                value={data.name}
                                onChange={(e) => setData('name', e.target.value)}
                                required
                                autoComplete="name"
                                placeholder="请输入您的姓名"
                            />

                            <InputError className="mt-2" message={errors.name} />
                        </div>

                        <div className="grid gap-2">
                            <Label htmlFor="username">用户名</Label>

                            <Input
                                id="username"
                                type="text"
                                className="mt-1 block w-full"
                                value={auth.user.username}
                                disabled
                                autoComplete="username"
                                placeholder="用户名不可更改"
                            />

                            <p className="text-sm text-muted-foreground">用户名不可更改</p>
                        </div>



                        <div className="flex items-center gap-4">
                            <Button disabled={processing}>保存</Button>

                            <Transition
                                show={recentlySuccessful}
                                enter="transition ease-in-out"
                                enterFrom="opacity-0"
                                leave="transition ease-in-out"
                                leaveTo="opacity-0"
                            >
                                <p className="text-sm text-neutral-600">已保存</p>
                            </Transition>
                        </div>
                    </form>
                </div>


            </SettingsLayout>
        </AppLayout>
    );
}
