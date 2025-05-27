export default defineNuxtPlugin(() => {
  return {
    provide: {
      authService: {
        async login(token: string) {
          const { $cookies } = useNuxtApp();
          $cookies.setItem('token', token);
        },
      
        async logout() {
          const { $cookies } = useNuxtApp();
          const headers = new Headers({
            'Content-Type': 'application/json',
            Authorization: `Bearer ${$cookies.getItem('token')}`,
          });
      
          try {
            const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/logout', {
              method: 'POST',
              headers: headers,
            });
      
            if (response.ok) {
              const message = await response.json();
              console.log('Logout successful:', message);
            } else {
              console.error('Failed to logout:', response.statusText);
            }
          } catch (error) {
            console.error('Error:', error);
          }
      
          $cookies.removeItem('token');
        },
      
        isAuthenticated() {
          const { $cookies } = useNuxtApp();
          const token = $cookies.getItem('token');
          if (!token) {
            return false;
          }

          const parsedToken = this.parseJwt(token);
          const currentTime = Math.floor(Date.now() / 1000);

          if(process.client) {
            if (parsedToken.exp && parsedToken.exp > currentTime) {
              return true;
            } else {
              // this.logout();
              return false;
            }
          } else {
            return false;
          }
        },
      
        getToken() {
          const { $cookies } = useNuxtApp();
          return $cookies.getItem('token');
        },
      
        parseJwt(token: string) {
          try {
            const base64Url = token.split('.')[1];
            const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
            const jsonPayload = decodeURIComponent(
              window
                .atob(base64)
                .split('')
                .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
                .join('')
            );
      
            return JSON.parse(jsonPayload);
          } catch {
            this.logout();
            navigateTo('/account/login');
          }
        },

        getUserTypeFromJWT() : number {
          const token = this.getToken();
          if (!token) {
            return 0;
          }
          const parsedToken = this.parseJwt(token);
          if (parsedToken.User && parsedToken.User.usertype) {
            return parsedToken.User.usertype;
          } else {
            return 0;
          }
        },

        isDocente() : boolean {
          return this.getUserTypeFromJWT() === Number(useRuntimeConfig().public.USER_TYPE_DOCENTE);
        },

        isDiscente() : boolean {
          return this.getUserTypeFromJWT() === Number(useRuntimeConfig().public.USER_TYPE_DISCENTE);
        },
      }
    }
  }
});