<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Foundation\Auth\User as Authenticatable;
use Illuminate\Notifications\Notifiable;

class AdminUser extends Authenticatable
{
    /** @use HasFactory<\Database\Factories\AdminUserFactory> */
    use HasFactory, Notifiable;

    /**
     * 数据库表名
     *
     * @var string
     */
    protected $table = 'admin_users';

    /**
     * The attributes that are mass assignable.
     *
     * @var list<string>
     */
    protected $fillable = [
        'name',
        'username',
        'password',
    ];

    /**
     * The attributes that should be hidden for serialization.
     *
     * @var list<string>
     */
    protected $hidden = [
        'password',
        'remember_token',
    ];

    /**
     * Get the attributes that should be cast.
     *
     * @return array<string, string>
     */
    protected function casts(): array
    {
        return [
            'password' => 'hashed',
        ];
    }

    /**
     * 获取用于认证的用户名字段
     *
     * @return string
     */
    public function getAuthIdentifierName(): string
    {
        return 'username';
    }

    /**
     * 获取用于密码重置的字段名
     *
     * @return string
     */
    public function getEmailForPasswordReset(): string
    {
        return $this->username;
    }
}
