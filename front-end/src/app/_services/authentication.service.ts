import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { User } from '@/_models';

@Injectable({ providedIn: 'root' })
export class AuthenticationService {
    private currentUserSubject: BehaviorSubject<User>;
    public currentUser: Observable<User>;
    public errorString: string;
    public homeCount: number;
    public mapPosInfoArray: object[];
    constructor(private http: HttpClient) {
        this.currentUserSubject = new BehaviorSubject<User>(JSON.parse(localStorage.getItem('currentUser')));
        this.currentUser = this.currentUserSubject.asObservable();

    }

    public get currentUserValue(): User {
        return this.currentUserSubject.value;
    }

    login(Email, HashedPassword) {
        return this.http.post<any>(`${config.apiUrl}/login`, { Email, HashedPassword })
            .pipe(map(user => {


                // store user details and jwt token in local storage to keep user logged in between page refreshes
                localStorage.setItem('currentUser', JSON.stringify(user.user));
                localStorage.setItem('currentInfo', JSON.stringify(user.mapposinfo));
                // this.currentUserSubject.next(user);
                this.currentUser = user.user
                return user;
            }));
    }
    getmapinfo() {

        return this.http.post<any>(`${config.apiUrl}/home`, {})
            .pipe(map(user => {


                this.mapPosInfoArray = user.mapposinfo;
                localStorage.setItem('currentInfo', JSON.stringify(user.mapposinfo));
                //localStorage.setItem('currentInfo',JSON.stringify(user.mapposinfo));
                // this.currentUserSubject.next(user);
                return user;
            }));
    }
    logout() {
        // remove user from local storage and set current user to null
        localStorage.removeItem('currentUser');
        this.currentUserSubject.next(null);
    }
}