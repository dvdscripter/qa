<template>
    <b-navbar variant="info" type="dark" toggleable>
        
        <b-navbar-nav>
            <b-nav-item :to="{name: 'Questions'}">
                Home
            </b-nav-item>
        </b-navbar-nav>

        <b-navbar-nav class="ml-auto">
            <b-nav-item-dropdown v-if="user" right>
                <template slot="button-content">
                    {{user}}
                </template>
                <b-dropdown-item :to="{name: 'Profile', params: {email: user}}">Profile</b-dropdown-item>
                <b-dropdown-item :to="{name: 'Logout'}">Signout</b-dropdown-item>
            </b-nav-item-dropdown>
            <b-nav-item v-if="!user" :to="{name: 'Login'}">
                Login
            </b-nav-item>
        </b-navbar-nav>
    </b-navbar>
</template>

<script>
export default {
  name: "NavBar",
  props: ["user"]
};
</script>