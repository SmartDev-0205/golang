import { Component, AfterViewInit, ViewChild, ElementRef, OnInit } from '@angular/core';
import { first } from 'rxjs/operators';

import { User } from '@/_models';
import { UserService, AuthenticationService } from '@/_services';
import { DH_NOT_SUITABLE_GENERATOR } from 'constants';
import { Router } from '@angular/router';
@Component({ templateUrl: 'home.component.html' })
export class HomeComponent implements AfterViewInit {
    @ViewChild('mapContainer', { static: false }) gmap: ElementRef;
    map: google.maps.Map
    constructor(
        private authenticationService: AuthenticationService,
        private router: Router,
        private userService: UserService
    ) {

        this.currentUser = this.authenticationService.currentUserValue;

    }

    //lat = 40.730610;
    //lng = -73.935242;
    coordinates = new google.maps.LatLng(this.authenticationService.mapPosInfoArray[0]["latitude"], this.authenticationService.mapPosInfoArray[0]["longitude"]);
    mapOptions: google.maps.MapOptions = {
        center: this.coordinates,
        zoom: 8,
    };

    icon = "./src/app/image/arrow.png";
    marker = new google.maps.Marker

    currentUser: User;
    users = [];
    //loading = false;

    ngAfterViewInit() {
        this.mapInitializer();
        setInterval(() => {
            this.getregularinfo();
        }, 1000);
    }
    logout() {
        //this.authenticationService.logout();
        localStorage.removeItem('currentUser')
        localStorage.removeItem('currentInfo')
        this.router.navigate(['/login']);

        window.location.reload();
    }
    getregularinfo() {
        this.authenticationService.getmapinfo().subscribe(
            data => {
                this.authenticationService.mapPosInfoArray = JSON.parse(localStorage.getItem('currentInfo'))

            },
            error => {
                alert("cannot get info from server")
            }
        );
        this.mapMarkerInitializer()
        this.displayMapTrack();
    }


    mapMarkerInitializer() {
        for (let index = 0; index < this.authenticationService.mapPosInfoArray.length; index++) {
            let angle = 0;
            if (index != 0) {
                angle = this.angleFromCoordinate(this.authenticationService.mapPosInfoArray[index - 1]["latitude"], this.authenticationService.mapPosInfoArray[index - 1]["longitude"],
                    this.authenticationService.mapPosInfoArray[index]["latitude"], this.authenticationService.mapPosInfoArray[index]["longitude"]);
            }
            this.marker = new google.maps.Marker({
                position: new google.maps.LatLng(this.authenticationService.mapPosInfoArray[index]["latitude"], this.authenticationService.mapPosInfoArray[index]["longitude"]),
                icon: {
                    path: google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
                    scale: 2,
                    rotation: angle
                },
                map: this.map,


            });
            this.marker.setMap(this.map);

        }
    }
    angleFromCoordinate(lat1, long1, lat2, long2) {
        let dLon = (long2 - long1);
        let y = Math.sin(dLon) * Math.cos(lat2);
        let x = Math.cos(lat1) * Math.sin(lat2) - Math.sin(lat1)
            * Math.cos(lat2) * Math.cos(dLon);
        let brng = Math.atan2(y, x);
        let pi = Math.PI;
        brng = brng * 180 / pi;
        brng = (brng + 360) % 360;
        brng = 360 - brng; // count degrees counter-clockwise - remove to make clockwise

        return brng;
    }

    mapInitializer() {

        this.map = new google.maps.Map(this.gmap.nativeElement, this.mapOptions);


    }
    displayMapTrack() {
        const flightPlanCoordinates = [];

        for (let index = 0; index < this.authenticationService.mapPosInfoArray.length; index++) {
            flightPlanCoordinates.push(new google.maps.LatLng(this.authenticationService.mapPosInfoArray[index]["latitude"], this.authenticationService.mapPosInfoArray[index]["longitude"]));

        }


        const flightPath = new google.maps.Polyline({
            path: flightPlanCoordinates,
            geodesic: true,
            strokeColor: "#000000",
            strokeOpacity: 1.0,
            strokeWeight: 2,
        });

        flightPath.setMap(this.map);
    }




}