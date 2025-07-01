<?php

use App\Admin\Controllers\Auth\AuthenticatedSessionController;
use App\Admin\Controllers\Auth\ConfirmablePasswordController;


use App\Admin\Controllers\Settings\PasswordController;
use App\Admin\Controllers\Settings\ProfileController;
use Inertia\Inertia;

// 根路由重定向到管理后台
Route::get('/', fn () => redirect('admin'))->name('home');

// 游客可访问的路由组
Route::group(['middleware' => 'guest', 'prefix' => 'admin'], function () {
    // 登录页面
    Route::get('login', [AuthenticatedSessionController::class, 'create']) ->name('login');
    
    // 登录处理
    Route::post('login', [AuthenticatedSessionController::class, 'store']);
});

// 需要认证的路由组
Route::group(['middleware' => 'auth', 'prefix' => 'admin'], function () {
    // 仪表盘
    Route::get('dashboard', fn () => Inertia::render('dashboard'))->name('dashboard');

    // 密码确认
    Route::get('confirm-password', [ConfirmablePasswordController::class, 'show'])
        ->name('password.confirm');
    Route::post('confirm-password', [ConfirmablePasswordController::class, 'store']);

    // 退出登录
    Route::post('logout', [AuthenticatedSessionController::class, 'destroy'])
        ->name('logout');

    // 设置页面路由
    Route::redirect('settings', 'settings/profile');

    // 个人资料设置
    Route::get('settings/profile', [ProfileController::class, 'edit'])->name('profile.edit');
    Route::patch('settings/profile', [ProfileController::class, 'update'])->name('profile.update');

    // 密码设置
    Route::get('settings/password', [PasswordController::class, 'edit'])->name('password.edit');
    Route::put('settings/password', [PasswordController::class, 'update'])->name('password.update');

    // 外观设置
    Route::get('settings/appearance', fn () => Inertia::render('settings/appearance'))->name('appearance');
});