// This example requires the Visualization library. Include the libraries=visualization
// parameter when you first load the API. For example:
// <script src="https://maps.googleapis.com/maps/api/js?key=YOUR_API_KEY&libraries=visualization">

var map, heatmap;

var mapData;

var gradient = [
"rgba(102, 255, 255, 1)",
"rgba(102, 255, 128, 1)",
"rgba(102, 255, 64, 1)",
"rgba(102, 255, 0, 1)",
"rgba(147, 255, 0, 1)",
"rgba(193, 255, 0, 1)",
"rgba(238, 255, 0, 1)",
"rgba(244, 227, 0, 1)",
"rgba(249, 198, 0, 1)",
"rgba(255, 170, 0, 1)",
"rgba(255, 113, 0, 1)",
"rgba(255, 57, 0, 1)",
"rgba(255, 0, 0, 1)"
];

document.getElementById("legend").style.background = 
"linear-gradient(to bottom, " + gradient + ")";
//document.getElementById("min")

function toggleCityDropdown() {
    document.getElementById("cityDropdown").classList.toggle("show");
}
function toggleTypeDropdown() {
    document.getElementById("typeDropdown").classList.toggle("show");
}

function filterFunction() {
    var input, filter, ul, li, a, i;
    input = document.getElementById("myInput");
    filter = input.value.toUpperCase();
    div = document.getElementById("cityDropdown");
    a = div.getElementsByTagName("a");
    for (i = 0; i < a.length; i++) {
        if (a[i].innerHTML.toUpperCase().indexOf(filter) > -1) {
            a[i].style.display = "";
        } else {
            a[i].style.display = "none";
        }
    }
}

function initMap() {
  mapData = getPolPoints();//new google.maps.MVCArray([]);
  map = new google.maps.Map(document.getElementById('map'), {
    zoom: 14,
    center: {lat: 53.2245534, lng: 6.571995},
    mapTypeId: 'roadmap'
  });

  heatmap = new google.maps.visualization.HeatmapLayer({
    data: mapData
  });
  heatmap.setMap(map);
  
  heatmap.set('radius', 144);
  heatmap.set('gradient', gradient);
  
  var opt = { minZoom: 8, maxZoom: 16 };
  map.setOptions(opt);
  
  heatmap.set('maxIntensity', 100);
  
  map.addListener('zoom_changed', function() {
    zoomLevel = map.getZoom();
    heatmap.set('radius', 9 * Math.pow(2,zoomLevel - 10));
  });
}

function toggleHeatmapPol() {
  heatmapPol.setMap(heatmapPol.getMap() ? null : map);
}
function toggleHeatmapTree() {
  heatmapTree.setMap(heatmapTree.getMap() ? null : map);
}

function Get(yourUrl){
    var Httpreq = new XMLHttpRequest(); // a new request
    Httpreq.open("GET",yourUrl,false);
    Httpreq.send(null);
    return Httpreq;
}

function getPolPoints() {
    var response = Get('http://localhost:8081/api/get/measurements_with_location?sensor_type=CO2')
    var maps = [];
    if (response.status != 200) {
        return maps
    }
    var jsonData = JSON.parse(response.responseText);

    for (var i = 0; i < jsonData.length; i++) {
        var row = jsonData[i];

        maps.push({location: new google.maps.LatLng(row.Longtitude, row.latitude), weight: row.measurement_value})
    }


    return maps;
}

function getGradePoints() {
    var response = Get('http://localhost:8081/api/get/measurements_with_location?sensor_type=grade')
    var maps = [];
    if (response.status != 200) {
        return maps
    }
    var jsonData = JSON.parse(response.responseText);

    for (var i = 0; i < jsonData.length; i++) {
        var row = jsonData[i];

        maps.push({location: new google.maps.LatLng(row.Longtitude, row.latitude), weight: row.measurement_value})
    }


    return maps;
}
