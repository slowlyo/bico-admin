<?php

namespace App\Admin\Requests;

use App\Models\AdminUser;
use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class AdminUserStoreRequest extends FormRequest
{
    /**
     * 确定用户是否有权限进行此请求
     */
    public function authorize(): bool
    {
        return true;
    }

    /**
     * 获取应用于请求的验证规则
     *
     * @return array<string, ValidationRule|array<mixed>|string>
     */
    public function rules(): array
    {
        return [
            'name' => ['required', 'string', 'max:255'],
            'username' => [
                'required',
                'string',
                'max:255',
                'alpha_dash',
                Rule::unique(AdminUser::class, 'username'),
            ],
            'password' => [
                'required',
                'string',
                'min:6',
                'confirmed',
            ],
        ];
    }

    /**
     * 获取验证错误的自定义属性名称
     *
     * @return array<string, string>
     */
    public function attributes(): array
    {
        return [
            'name' => '姓名',
            'username' => '用户名',
            'password' => '密码',
        ];
    }

    /**
     * 获取验证错误的自定义消息
     *
     * @return array<string, string>
     */
    public function messages(): array
    {
        return [
            'username.alpha_dash' => '用户名只能包含字母、数字、破折号和下划线',
            'username.unique' => '该用户名已被使用',
            'password.min' => '密码至少需要6个字符',
            'password.confirmed' => '密码确认不匹配',
        ];
    }
}
