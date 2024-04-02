// Client side logic
// Author: Justin Forseth

// Default starting position
CARROLL_COORDS = [46.60085, -112.03872];

// Active marker on screen
var marker;

// Flag to prevent hover effect from activating when creating a new class
var hoverLocked = false;

// Initialize map
// This assumes Leaflet is loaded correctly
var map = L.map("map").setView(CARROLL_COORDS, 17);
L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 40,
  attribution:
    '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
}).addTo(map);


map.on("click", addRoom);

// When a card is hovered, zoom in on position 
// specified by data-lat and data-long
function hover(event) {
  // Bail if hover effect is disabled
  if (hoverLocked) {
    return;
  }
  
  // Get position data
  const lat = event.target.dataset.lat;
  const long = event.target.dataset.long;
  
  // Remove old marker if there is one
  if (marker != null) {
    marker.remove();
  }
  // Create new marker
  marker = L.marker([lat, long]).addTo(map);
  
  // Zoom in
  map.flyTo([lat, long], 20);
}

// Go to default position
function zoomDefault(){
  map.flyTo(CARROLL_COORDS, 17);
}

// Zoom back out when card is unhovered
function leave() {
  // Bail if hover effect is disabled
  if (hoverLocked) {
    return;
  }
  // Remove old marker if there is one
  if (marker != null) {
    marker.remove();
  }
  zoomDefault();
}

// Zoom in on selected building, 
// and populate roomsSelect dropdown
async function selectBuilding(event) {
  // Disable card hovering, we care about the new class,
  // not existing ones
  hoverLocked = true;

  // Get the id of the selected building
  const select = event.target;
  const selectedBuilding = select.options[select.selectedIndex];
  const buildingId = selectedBuilding.value;
  
  // Zoom in on the selected building
  map.flyTo([selectedBuilding.dataset.lat, selectedBuilding.dataset.long], 20);

  // Get the list of rooms from the server
  const response = await fetch(`/api/rooms/${buildingId}`);
  const rooms = await response.json();

  // Update the roomsSelect dropdown with the fetched room data
  const roomsSelect = document.getElementById("roomsSelect");
  roomsSelect.innerHTML = ""; // Clear existing options

  rooms.forEach((room) => {
    const option = document.createElement("option");
    option.value = room.ID; 
    option.textContent = room.Name; 
    option.dataset.lat = room.Latitude;
    option.dataset.long = room.Longitude;
    roomsSelect.appendChild(option);
  });
}

// When a room is selected from the dropdown, mark it and zoom in
async function selectRoom(event) {
  // Disable card hovering
  hoverLocked = true;
  
  // Get selected room
  const select = event.target;
  const selectedRoom = select.options[select.selectedIndex];

  // Remove old marker if there is one
  if (marker != null) {
    marker.remove();
  }
  // Add new marker
  marker = L.marker([
    selectedRoom.dataset.lat,
    selectedRoom.dataset.long,
  ]).addTo(map);
  // Zoom in
  map.flyTo([selectedRoom.dataset.lat, selectedRoom.dataset.long], 19);
}

// Create a new room when the map is clicked
async function addRoom(e) {
  // Get room select (so we can set the select ro
  
  // Get building select (so we know what building the room is in
  const buildingSelect = document.getElementById("buildingSelect"); // Add this line

  // Prompt user for room name
  const roomName = prompt("Enter room name:");
  if (!roomName) {
    // User canceled the prompt
    return;
  }

  // Create an object with room data
  const roomData = {
    Name: roomName,
    Latitude: e.latlng.lat,
    Longitude: e.latlng.lng,
    BuildingID: Number(buildingSelect.value),
  };

  // Make API call to submit room data
  const response = await fetch("/api/rooms/add-room", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(roomData),
  });

  if (response.ok) {
    // API call successful, parse the response
    const { roomId } = await response.json();

    // Add room as an option to the select element
    const option = document.createElement("option");
    option.value = roomId; // Assuming the ID is returned from the API
    option.textContent = roomName;
    option.dataset.lat = roomData.Latitude;
    option.dataset.long = roomData.Longitude;
    
    const roomsSelect = document.getElementById("roomsSelect");
    roomsSelect.appendChild(option);
    // Reenable hovering
    hoverLocked = false;
  }
}
