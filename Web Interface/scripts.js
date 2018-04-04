// This example requires the Visualization library. Include the libraries=visualization
// parameter when you first load the API. For example:
// <script src="https://maps.googleapis.com/maps/api/js?key=YOUR_API_KEY&libraries=visualization">

var map, heatmap;

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

function toggleDropdown() {
    document.getElementById("cityDropdown").classList.toggle("show");
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
  map = new google.maps.Map(document.getElementById('map'), {
    zoom: 14,
    center: {lat: 53.2245534, lng: 6.571995},
    mapTypeId: 'roadmap'
  });

  heatmap = new google.maps.visualization.HeatmapLayer({
    data: getPoints(),
    map: map
  });
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

function toggleHeatmap() {
  heatmap.setMap(heatmap.getMap() ? null : map);
}

function changeOpacity() {
  heatmap.set('opacity', heatmap.get('opacity') ? null : 0.2);
}

function getPoints() {
  return [
    {location: new google.maps.LatLng(53.203246, 6.564907), weight: 58.1},
    {location: new google.maps.LatLng(53.205686, 6.57557), weight: 62.36},
    {location: new google.maps.LatLng(53.210628, 6.586484), weight: 62.52},
    {location: new google.maps.LatLng(53.213608, 6.597814), weight: 60.45},
    {location: new google.maps.LatLng(53.217179, 6.609401), weight: 56.89},
    {location: new google.maps.LatLng(53.220467, 6.615795), weight: 64.26},
    {location: new google.maps.LatLng(53.224166, 6.613649), weight: 61.87},
    {location: new google.maps.LatLng(53.229355, 6.61013), weight: 62.48},
    {location: new google.maps.LatLng(53.23686, 6.593625), weight: 62.3},
    {location: new google.maps.LatLng(53.241534, 6.585857), weight: 64.87},
    {location: new google.maps.LatLng(53.24631, 6.581737), weight: 56.29},
    {location: new google.maps.LatLng(53.248903, 6.576501), weight: 60.49},
    {location: new google.maps.LatLng(53.246745, 6.572167), weight: 64.88},
    {location: new google.maps.LatLng(53.240735, 6.568562), weight: 64.92},
    {location: new google.maps.LatLng(53.238346, 6.561138), weight: 61.57},
    {location: new google.maps.LatLng(53.237575, 6.548564), weight: 64.66},
    {location: new google.maps.LatLng(53.235956, 6.539037), weight: 57.28},
    {location: new google.maps.LatLng(53.234491, 6.528265), weight: 57.5},
    {location: new google.maps.LatLng(53.230483, 6.531441), weight: 64.86},
    {location: new google.maps.LatLng(53.226757, 6.534488), weight: 55.95},
    {location: new google.maps.LatLng(53.222306, 6.538189), weight: 58.64},
    {location: new google.maps.LatLng(53.213137, 6.541998), weight: 55.8},
    {location: new google.maps.LatLng(53.207585, 6.548779), weight: 58.37},
    {location: new google.maps.LatLng(53.202623, 6.552298), weight: 59.74},
    {location: new google.maps.LatLng(53.234868, 6.603167), weight: 57.77},
    {location: new google.maps.LatLng(53.21787, 6.539329), weight: 61.25},
    {location: new google.maps.LatLng(53.202519, 6.559937), weight: 55.42},
    {location: new google.maps.LatLng(53.208011, 6.580399), weight: 61.57},
    {location: new google.maps.LatLng(53.21297, 6.592201), weight: 56.54},
    {location: new google.maps.LatLng(53.215359, 6.604131), weight: 63.95},
    {location: new google.maps.LatLng(53.204229, 6.569756), weight: 63.84},
    {location: new google.maps.LatLng(53.20586, 6.542687), weight: 0.65},
    {location: new google.maps.LatLng(53.203571, 6.540241), weight: 1.98},
    {location: new google.maps.LatLng(53.199637, 6.537838), weight: 3.9},
    {location: new google.maps.LatLng(53.218449, 6.532525), weight: 2.24},
    {location: new google.maps.LatLng(53.224975, 6.516904), weight: 3.79},
    {location: new google.maps.LatLng(53.241645, 6.543769), weight: 2.82},
    {location: new google.maps.LatLng(53.240437, 6.551451), weight: 3.36},
    {location: new google.maps.LatLng(53.235017, 6.570892), weight: 3.29},
    {location: new google.maps.LatLng(53.228697, 6.582393), weight: 0.99},
    {location: new google.maps.LatLng(53.228131, 6.548704), weight: 0.83},
    {location: new google.maps.LatLng(53.226409, 6.558274), weight: 1.82},
    {location: new google.maps.LatLng(53.224301, 6.555227), weight: 2.3},
    {location: new google.maps.LatLng(53.221371, 6.554583), weight: 3.95},
    {location: new google.maps.LatLng(53.232252, 6.593819), weight: 3.61},
    {location: new google.maps.LatLng(53.192011, 6.540175), weight: 3.22},
    {location: new google.maps.LatLng(53.195422, 6.547616), weight: 1.83},
    {location: new google.maps.LatLng(53.22675, 6.586754), weight: 1.29},
    {location: new google.maps.LatLng(53.221456, 6.553623), weight: 2.78},
    {location: new google.maps.LatLng(53.229756, 6.546128), weight: 22.32},
    {location: new google.maps.LatLng(53.210355, 6.562007), weight: 44.57},
    {location: new google.maps.LatLng(53.232671, 6.557447), weight: 33.12},
    {location: new google.maps.LatLng(53.234489, 6.580317), weight: 19.88},
    {location: new google.maps.LatLng(53.226191, 6.59817), weight: 22.53},
    {location: new google.maps.LatLng(53.220229, 6.590316), weight: 23.44},
    {location: new google.maps.LatLng(53.214806, 6.557443), weight: 39.66},
    {location: new google.maps.LatLng(53.213931, 6.573579), weight: 41.71},
    {location: new google.maps.LatLng(53.227369, 6.570789), weight: 20.89},
    {location: new google.maps.LatLng(53.222538, 6.56654), weight: 23.19},
    {location: new google.maps.LatLng(53.221972, 6.578428), weight: 32.66}
  ];
}
