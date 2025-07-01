// Components
import { Head, useForm } from '@inertiajs/react';
import { LoaderCircle } from 'lucide-react';
import { FormEventHandler } from 'react';

import TextLink from '@/components/text-link';
import { Button } from '@/components/ui/button';
import AuthLayout from '@/layouts/auth-layout';

export default function VerifyEmail({ status }: { status?: string }) {
    const { post, processing } = useForm({});

    const submit: FormEventHandler = (e) => {
        e.preventDefault();

        post(route('verification.send'));
    };

    return (
        <AuthLayout title="验证邮箱" description="请点击我们刚刚发送到您邮箱的链接来验证您的邮箱地址。">
            <Head title="邮箱验证" />

            {status === 'verification-link-sent' && (
                <div className="mb-4 text-center text-sm font-medium text-green-600">
                    新的验证链接已发送到您注册时提供的邮箱地址。
                </div>
            )}

            <form onSubmit={submit} className="space-y-6 text-center">
                <Button disabled={processing} variant="secondary">
                    {processing && <LoaderCircle className="h-4 w-4 animate-spin" />}
                    重新发送验证邮件
                </Button>

                <TextLink href={route('logout')} method="post" className="mx-auto block text-sm">
                    退出登录
                </TextLink>
            </form>
        </AuthLayout>
    );
}
