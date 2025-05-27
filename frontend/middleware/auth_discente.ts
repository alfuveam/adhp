export default defineNuxtRouteMiddleware((to, from) => {
    if(process.client) {
        const { $authService } = useNuxtApp();
        if (!$authService.isAuthenticated()) {
            return navigateTo('/account/login');
        }
        // if (!$authService.isDiscente()) {
        //     return navigateTo('/account/login');
        // }
        navigateTo(to.path);
    }
})