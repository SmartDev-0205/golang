import { Injectable } from '@angular/core';
import { Router, CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';
import { AuthenticationService } from '@/_services';

@Injectable({ providedIn: 'root' })
export class AuthGuard implements CanActivate {
    constructor(
        private router: Router,
        private authenticationService: AuthenticationService,
        private http: HttpClient
    ) { }

    canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot) {

        this.authenticationService.mapPosInfoArray = JSON.parse(localStorage.getItem('currentInfo'))
        const currentUser = JSON.parse(localStorage.getItem('currentUser'));
        if (currentUser) {

            return true;
        }

        // not logged in so redirect to login page with the return url

        this.router.navigate(['/login']);
        return false;
    }
}