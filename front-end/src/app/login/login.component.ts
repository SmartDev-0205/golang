import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { first } from 'rxjs/operators';

import { AlertService, AuthenticationService } from '@/_services';

@Component({ templateUrl: 'login.component.html' })
export class LoginComponent implements OnInit {
    loginForm: FormGroup;
    loading = false;
    submitted = false;
    returnUrl: string;

    constructor(
        private formBuilder: FormBuilder,
        private route: ActivatedRoute,
        private router: Router,
        private authenticationService: AuthenticationService,
        private alertService: AlertService
    ) {
        
        // redirect to home if already logged in
        if (this.authenticationService.currentUserValue) {
            this.router.navigate(['/']);
        }
    }

    ngOnInit() {
        this.loginForm = this.formBuilder.group({
            Email: ['', Validators.required],
            HashedPassword: ['', Validators.required]
        });

        // get return url from route parameters or default to '/'
       // this.returnUrl = this.route.snapshot.queryParams['returnUrl'] || '/';

    }

    // convenience getter for easy access to form fields
    get f() { return this.loginForm.controls; }

    onSubmit() {
        this.submitted = true;

        // reset alerts on submit
        this.alertService.clear();

        // stop here if form is invalid
        if (this.loginForm.invalid) {
            return;
        }

        this.loading = false;
        this.authenticationService.login(this.f.Email.value, this.f.HashedPassword.value)
            .pipe(first())
            .subscribe(
                data => {
                    console.log(data)
                    this.authenticationService.errorString=data.error;
                    if (data.error==""){
                        this.authenticationService.mapPosInfoArray=data.mapposinfo
                        console.log(this.authenticationService.mapPosInfoArray)
                        //this.authenticationService.homeCount=0
                        this.router.navigate(['']);
                    }else{
                        this.alertService.error(data.error);
                    }
                    
                },
                error => {
                    console.log("error:",error)
                    this.alertService.error(error);
                    this.loading = false;
                });
    }
}
