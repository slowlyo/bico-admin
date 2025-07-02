<?php

namespace App\Admin\Controllers\System;

use App\Admin\Controllers\Controller;
use App\Models\AdminUser;
use App\Admin\Requests\AdminUserStoreRequest;
use App\Admin\Requests\AdminUserUpdateRequest;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;
use Inertia\Inertia;
use Inertia\Response;

class UserController extends Controller
{
    /**
     * 显示用户列表
     */
    public function index(Request $request): Response
    {
        $query = AdminUser::query();

        // 搜索功能
        if ($request->filled('search')) {
            $search = $request->get('search');
            $query->where(function ($q) use ($search) {
                $q->where('name', 'like', "%{$search}%")
                  ->orWhere('username', 'like', "%{$search}%");
            });
        }

        // 排序功能
        $sortField = $request->get('sort', 'created_at');
        $sortDirection = $request->get('direction', 'desc');
        
        $allowedSorts = ['name', 'username', 'created_at', 'updated_at'];
        if (in_array($sortField, $allowedSorts)) {
            $query->orderBy($sortField, $sortDirection);
        }

        $users = $query->paginate(15)->withQueryString();

        return Inertia::render('users/index', [
            'users' => $users,
            'filters' => [
                'search' => $request->get('search'),
                'sort' => $sortField,
                'direction' => $sortDirection,
            ],
        ]);
    }

    /**
     * 显示创建用户表单
     */
    public function create(): Response
    {
        return Inertia::render('users/create');
    }

    /**
     * 存储新用户
     */
    public function store(AdminUserStoreRequest $request): RedirectResponse
    {
        $validated = $request->validated();
        
        AdminUser::create([
            'name' => $validated['name'],
            'username' => $validated['username'],
            'password' => Hash::make($validated['password']),
        ]);

        return redirect()->route('admin.users.index')
            ->with('success', '用户创建成功');
    }

    /**
     * 显示用户详情
     */
    public function show(AdminUser $user): Response
    {
        return Inertia::render('users/show', [
            'user' => $user,
        ]);
    }

    /**
     * 显示编辑用户表单
     */
    public function edit(AdminUser $user): Response
    {
        return Inertia::render('users/edit', [
            'user' => $user,
        ]);
    }

    /**
     * 更新用户信息
     */
    public function update(AdminUserUpdateRequest $request, AdminUser $user): RedirectResponse
    {
        $validated = $request->validated();
        
        $updateData = [
            'name' => $validated['name'],
        ];

        // 如果提供了新密码，则更新密码
        if (!empty($validated['password'])) {
            $updateData['password'] = Hash::make($validated['password']);
        }

        $user->update($updateData);

        return redirect()->route('admin.users.index')
            ->with('success', '用户信息更新成功');
    }

    /**
     * 删除用户
     */
    public function destroy(AdminUser $user): RedirectResponse
    {
        // 防止用户删除自己的账户
        if ($user->id === auth('admin')->id()) {
            return redirect()->route('admin.users.index')
                ->with('error', '不能删除自己的账户');
        }

        $user->delete();

        return redirect()->route('admin.users.index')
            ->with('success', '用户删除成功');
    }
}
