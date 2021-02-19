import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { User } from '@/_models';

@Injectable({ providedIn: 'root' })
export class UserService {
    constructor(private http: HttpClient) { }
    
    getAll() {
        return this.http.get<User[]>(`${config.apiUrl}/users`);
    }

    register(user: User) {
        console.log("user:",user)
        return this.http.post(`${config.apiUrl}/register`, user);
    }

    delete(id: number) {
        return this.http.delete(`${config.apiUrl}/users/${id}`);
    }
}
