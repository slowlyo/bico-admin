<?php

use App\Admin\Controllers\Auth\AuthenticatedSessionController;
use App\Admin\Controllers\Auth\ConfirmablePasswordController;
use App\Admin\Controllers\Auth\EmailVerificationNotificationController;
use App\Admin\Controllers\Auth\EmailVerificationPromptController;
use App\Admin\Controllers\Auth\NewPasswordController;
use App\Admin\Controllers\Auth\PasswordResetLinkController;
use App\Admin\Controllers\Auth\RegisteredUserController;
use App\Admin\Controllers\Auth\VerifyEmailController;

use App\Admin\Controllers\Settings\PasswordController;
use App\Admin\Controllers\Settings\ProfileController;
use Inertia\Inertia;

Route::get('/', function () {
    return redirect('admin');
})->name('home');

Route::middleware(['auth', 'verified'])->get('admin', function () {
    return Inertia::render('dashboard');
})->name('admin');

Route::middleware(['auth', 'verified'])->group(function () {
    Route::get('admin/dashboard', fn () => Inertia::render('dashboard'))->name('dashboard');
});

Route::group(['middleware' => 'guest', 'prefix' => 'admin'], function () {
    Route::get('register', [RegisteredUserController::class, 'create'])
        ->name('register');

    Route::post('register', [RegisteredUserController::class, 'store']);

    Route::get('login', [AuthenticatedSessionController::class, 'create'])
        ->name('login');

    Route::post('login', [AuthenticatedSessionController::class, 'store']);

    Route::get('forgot-password', [PasswordResetLinkController::class, 'create'])
        ->name('password.request');

    Route::post('forgot-password', [PasswordResetLinkController::class, 'store'])
        ->name('password.email');

    Route::get('reset-password/{token}', [NewPasswordController::class, 'create'])
        ->name('password.reset');

    Route::post('reset-password', [NewPasswordController::class, 'store'])
        ->name('password.store');
});

Route::group(['middleware' => 'auth', 'prefix' => 'admin'], function () {
    Route::get('verify-email', EmailVerificationPromptController::class)
        ->name('verification.notice');

    Route::get('verify-email/{id}/{hash}', VerifyEmailController::class)
        ->middleware(['signed', 'throttle:6,1'])
        ->name('verification.verify');

    Route::post('email/verification-notification', [EmailVerificationNotificationController::class, 'store'])
        ->middleware('throttle:6,1')
        ->name('verification.send');

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