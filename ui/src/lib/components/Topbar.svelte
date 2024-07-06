<script>
	import IconMenu from '~icons/mdi/menu';
	import IconBell from '~icons/mdi/bell';
	import IconAccount from '~icons/mdi/account-circle';

	const { toggleDrawer } = $props();
</script>

<div class="navbar bg-base-100 z-50 sticky top-0">
	<div class="flex-1">
		<button class="btn btn-square btn-ghost lg:hidden" onclick={toggleDrawer}>
			<IconMenu class="w-6 h-6" />
		</button>
		<a href="/" class="sm:hidden btn btn-ghost normal-case text-2xl">Ekolod</a>
	</div>
	<div class="flex-none gap-2">
		<button type="button" class="btn btn-ghost btn-circle" aria-label="Notifications">
			<div class="indicator">
				<IconBell class="w-5 h-5" />
				<span class="badge badge-xs badge-primary indicator-item"></span>
			</div>
		</button>
		<div class="dropdown dropdown-end">
			<button
				type="button"
				class="btn btn-ghost btn-circle"
				aria-haspopup="true"
				aria-expanded="false"
			>
				<IconAccount class="w-8 h-8 text-base-content" />
			</button>
			<ul class="mt-3 p-2 shadow menu menu-sm dropdown-content bg-base-100 rounded-box w-52">
				<li><a href="/profile">Profile</a></li>
				<li><a href="/settings">Settings</a></li>
				<li><a href="/logout">Logout</a></li>
			</ul>
		</div>
	</div>
</div>
