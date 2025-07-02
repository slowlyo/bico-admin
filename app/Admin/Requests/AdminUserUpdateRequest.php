<?php

namespace App\Admin\Requests;

use App\Models\AdminUser;
use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class AdminUserUpdateRequest extends FormRequest
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
        $userId = $this->route('user')->id;

        return [
            'name' => ['required', 'string', 'max:255'],
            'password' => [
                'nullable',
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
            'password.min' => '密码至少需要6个字符',
            'password.confirmed' => '密码确认不匹配',
        ];
    }
}
