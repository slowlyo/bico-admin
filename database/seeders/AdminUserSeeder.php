<?php

namespace Database\Seeders;

use App\Models\AdminUser;
use Illuminate\Database\Seeder;

class AdminUserSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        // 创建默认管理员用户
        AdminUser::factory()->create([
            'name' => '管理员',
            'username' => 'admin',
            'password' => bcrypt('admin'),
        ]);
    }
}
