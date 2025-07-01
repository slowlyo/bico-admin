<?php

use App\Admin\Controllers\Auth\AuthenticatedSessionController;
use App\Admin\Controllers\Auth\ConfirmablePasswordController;


use App\Admin\Controllers\Settings\PasswordController;
use App\Admin\Controllers\Settings\ProfileController;
use Inertia\Inertia;

Route::get('/', function () {
    return redirect('admin');
})->name('home');

Route::middleware(['auth'])->get('admin', function () {
    return Inertia::render('dashboard');
})->name('admin');

Route::middleware(['auth'])->group(function () {
    Route::get('admin/dashboard', fn () => Inertia::render('dashboard'))->name('dashboard');
});

Route::group(['middleware' => 'guest', 'prefix' => 'admin'], function () {
    Route::get('login', [AuthenticatedSessionController::class, 'create'])
        ->name('login');

    Route::post('login', [AuthenticatedSessionController::class, 'store']);
});

Route::group(['middleware' => 'auth', 'prefix' => 'admin'], function () {
    Route::get('confirm-password', [ConfirmablePasswordController::class, 'show'])
        ->name('password.confirm');

    Route::post('confirm-password', [ConfirmablePasswordController::class, 'store']);

    Route::post('logout', [AuthenticatedSessionController::class, 'destroy'])
        ->name('logout');

    Route::redirect('settings', 'settings/profile');

    Route::get('settings/profile', [ProfileController::class, 'edit'])->name('profile.edit');
    Route::patch('settings/profile', [ProfileController::class, 'update'])->name('profile.update');
    Route::delete('settings/profile', [ProfileController::class, 'destroy'])->name('profile.destroy');

    Route::get('settings/password', [PasswordController::class, 'edit'])->name('password.edit');
    Route::put('settings/password', [PasswordController::class, 'update'])->name('password.update');

    Route::get('settings/appearance', function () {
        return Inertia::render('settings/appearance');
    })->name('appearance');
});