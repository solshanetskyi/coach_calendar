package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Meeting Booking - Coach Calendar</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }

        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 10px;
        }

        .header p {
            font-size: 1.1rem;
            opacity: 0.9;
        }

        .content {
            padding: 40px;
        }

        .calendar-grid {
            display: grid;
            grid-template-columns: repeat(7, 1fr);
            gap: 10px;
            margin-bottom: 30px;
        }

        .day-cell {
            padding: 20px 10px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            text-align: center;
            cursor: pointer;
            transition: all 0.3s ease;
            background: white;
            min-height: 80px;
            display: flex;
            flex-direction: column;
            justify-content: center;
        }

        .day-cell:hover {
            border-color: #667eea;
            background: #f0f4ff;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
        }

        .day-cell.selected {
            border-color: #667eea;
            background: #667eea;
            color: white;
        }

        .day-cell.no-slots {
            background: #f5f5f5;
            color: #999;
            cursor: not-allowed;
            opacity: 0.5;
        }

        .day-cell.no-slots:hover {
            transform: none;
            box-shadow: none;
            background: #f5f5f5;
            border-color: #e0e0e0;
        }

        .day-number {
            font-size: 1.5rem;
            font-weight: bold;
            margin-bottom: 5px;
        }

        .day-name {
            font-size: 0.8rem;
            opacity: 0.8;
            text-transform: uppercase;
        }

        .day-slots-count {
            font-size: 0.75rem;
            margin-top: 5px;
            opacity: 0.9;
        }

        .time-slots-panel {
            background: #f9f9f9;
            border-radius: 12px;
            padding: 30px;
            margin-bottom: 30px;
            display: none;
        }

        .time-slots-panel.active {
            display: block;
        }

        .time-slots-header {
            text-align: center;
            margin-bottom: 20px;
        }

        .time-slots-header h3 {
            color: #333;
            font-size: 1.3rem;
            margin-bottom: 5px;
        }

        .time-slots-header p {
            color: #666;
            font-size: 0.9rem;
        }

        .time-slots-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
            gap: 10px;
            margin-bottom: 20px;
        }

        .time-slot {
            padding: 15px 10px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            text-align: center;
            cursor: pointer;
            transition: all 0.3s ease;
            background: white;
            font-size: 1rem;
            font-weight: 600;
        }

        .time-slot:hover {
            border-color: #667eea;
            background: #f0f4ff;
            transform: translateY(-2px);
        }

        .time-slot.booked {
            background: #f5f5f5;
            color: #999;
            cursor: not-allowed;
            opacity: 0.6;
            text-decoration: line-through;
        }

        .time-slot.booked:hover {
            transform: none;
            border-color: #e0e0e0;
            background: #f5f5f5;
        }

        .back-btn {
            padding: 10px 20px;
            background: #e0e0e0;
            color: #666;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s;
            display: block;
            margin: 0 auto 20px;
        }

        .back-btn:hover {
            background: #d0d0d0;
        }

        .booking-form {
            max-width: 500px;
            margin: 30px auto;
            padding: 30px;
            background: #f9f9f9;
            border-radius: 12px;
            display: none;
        }

        .booking-form.active {
            display: block;
        }

        .form-group {
            margin-bottom: 20px;
        }

        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
        }

        .form-group input {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            font-size: 1rem;
            transition: border-color 0.3s;
        }

        .form-group input:focus {
            outline: none;
            border-color: #667eea;
        }

        .selected-slot-info {
            background: #e8eaf6;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
            text-align: center;
        }

        .selected-slot-info strong {
            color: #667eea;
            font-size: 1.1rem;
        }

        .btn {
            padding: 14px 32px;
            border: none;
            border-radius: 8px;
            font-size: 1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
            width: 100%;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }

        .btn-secondary {
            background: #e0e0e0;
            color: #666;
            margin-top: 10px;
        }

        .btn-secondary:hover {
            background: #d0d0d0;
        }

        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
        }

        .month-navigation {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 30px;
            padding: 20px;
            background: #f9f9f9;
            border-radius: 12px;
        }

        .month-title {
            font-size: 1.5rem;
            font-weight: 600;
            color: #333;
        }

        .nav-btn {
            padding: 10px 20px;
            background: #667eea;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s;
        }

        .nav-btn:hover {
            background: #5568d3;
            transform: translateY(-2px);
        }

        .nav-btn:disabled {
            background: #ccc;
            cursor: not-allowed;
            transform: none;
        }

        .message {
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
            text-align: center;
            display: none;
        }

        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        .message.active {
            display: block;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Meeting Booking</h1>
            <p>Select a day, then choose your preferred time slot</p>
        </div>

        <div class="content">
            <div id="message" class="message"></div>

            <div id="loading" class="loading">
                Loading available slots...
            </div>

            <div id="monthNavigation" class="month-navigation" style="display: none;">
                <button class="nav-btn" id="prevMonth" onclick="changeMonth(-1)">← Previous</button>
                <div class="month-title" id="currentMonth"></div>
                <button class="nav-btn" id="nextMonth" onclick="changeMonth(1)">Next →</button>
            </div>

            <div id="calendar" class="calendar-grid"></div>

            <div id="timeSlotsPanel" class="time-slots-panel">
                <button class="back-btn" onclick="backToCalendar()">← Back to Calendar</button>
                <div class="time-slots-header">
                    <h3 id="selectedDateTitle"></h3>
                    <p>Choose an available time slot</p>
                </div>
                <div id="timeSlotsGrid" class="time-slots-grid"></div>
            </div>

            <div id="bookingForm" class="booking-form">
                <div class="selected-slot-info">
                    <div>Selected Time:</div>
                    <strong id="selectedSlotDisplay"></strong>
                </div>

                <div class="form-group">
                    <label for="name">Your Name</label>
                    <input type="text" id="name" placeholder="Enter your full name" required>
                </div>

                <div class="form-group">
                    <label for="email">Your Email</label>
                    <input type="email" id="email" placeholder="your.email@example.com" required>
                </div>

                <button class="btn btn-primary" onclick="confirmBooking()">Confirm Booking</button>
                <button class="btn btn-secondary" onclick="cancelBooking()">Cancel</button>
            </div>
        </div>
    </div>

    <script>
        let selectedSlot = null;
        let selectedDay = null;
        let allSlots = [];
        let currentMonthIndex = 0;
        let availableMonths = [];

        function formatDateTime(isoString) {
            const date = new Date(isoString);
            const options = {
                weekday: 'long',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            };
            return date.toLocaleString('en-US', options);
        }

        function formatDateLong(isoString) {
            const date = new Date(isoString);
            return date.toLocaleDateString('en-US', {
                weekday: 'long',
                month: 'long',
                day: 'numeric',
                year: 'numeric'
            });
        }

        function formatTime(isoString) {
            const date = new Date(isoString);
            return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: true });
        }

        function formatMonthYear(year, month) {
            const date = new Date(year, month);
            return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
        }

        function getDayKey(dateString) {
            const date = new Date(dateString);
            return date.getFullYear() + '-' + (date.getMonth() + 1) + '-' + date.getDate();
        }

        function groupSlotsByMonth(slots) {
            const monthMap = new Map();

            slots.forEach(slot => {
                const date = new Date(slot.slot_time);
                const year = date.getFullYear();
                const month = date.getMonth();
                const key = year + '-' + month;

                if (!monthMap.has(key)) {
                    monthMap.set(key, {
                        year: year,
                        month: month,
                        days: new Map()
                    });
                }

                const dayKey = getDayKey(slot.slot_time);
                const monthData = monthMap.get(key);

                if (!monthData.days.has(dayKey)) {
                    monthData.days.set(dayKey, {
                        date: new Date(date.getFullYear(), date.getMonth(), date.getDate()),
                        slots: []
                    });
                }

                monthData.days.get(dayKey).slots.push(slot);
            });

            return Array.from(monthMap.values()).sort((a, b) => {
                if (a.year !== b.year) return a.year - b.year;
                return a.month - b.month;
            });
        }

        function getMondayBasedWeekday(date) {
            // Convert Sunday (0) to 6, Monday (1) to 0, Tuesday (2) to 1, etc.
            const day = date.getDay();
            return day === 0 ? 6 : day - 1;
        }

        function changeMonth(direction) {
            currentMonthIndex += direction;
            if (currentMonthIndex < 0) currentMonthIndex = 0;
            if (currentMonthIndex >= availableMonths.length) {
                currentMonthIndex = availableMonths.length - 1;
            }
            renderCalendar();
        }

        function updateMonthNavigation() {
            if (availableMonths.length === 0) return;

            const monthData = availableMonths[currentMonthIndex];
            document.getElementById('currentMonth').textContent = formatMonthYear(monthData.year, monthData.month);
            document.getElementById('prevMonth').disabled = currentMonthIndex === 0;
            document.getElementById('nextMonth').disabled = currentMonthIndex === availableMonths.length - 1;
        }

        async function loadSlots() {
            try {
                const response = await fetch('/api/slots');
                if (!response.ok) {
                    throw new Error('Failed to load slots');
                }

                allSlots = await response.json();
                availableMonths = groupSlotsByMonth(allSlots);
                currentMonthIndex = 0;

                document.getElementById('loading').style.display = 'none';
                document.getElementById('monthNavigation').style.display = 'flex';

                renderCalendar();
            } catch (error) {
                console.error('Error loading slots:', error);
                showMessage('Failed to load available slots. Please refresh the page.', 'error');
                document.getElementById('loading').style.display = 'none';
            }
        }

        function renderCalendar() {
            const calendar = document.getElementById('calendar');
            calendar.innerHTML = '';

            if (availableMonths.length === 0) {
                calendar.innerHTML = '<div style="text-align: center; padding: 40px; color: #666;">No available slots found.</div>';
                return;
            }

            const monthData = availableMonths[currentMonthIndex];
            const days = Array.from(monthData.days.values()).sort((a, b) => {
                // First sort by date
                const dateDiff = a.date - b.date;
                if (dateDiff !== 0) return dateDiff;
                return 0;
            });

            // Group days by weeks (starting Monday)
            const weeks = [];
            let currentWeek = [];

            days.forEach((dayData, index) => {
                const weekday = getMondayBasedWeekday(dayData.date);

                // If this is Monday (0) and we have days in current week, start a new week
                if (weekday === 0 && currentWeek.length > 0) {
                    weeks.push(currentWeek);
                    currentWeek = [];
                }

                // Add empty cells at the start of the first week if it doesn't start on Monday
                if (index === 0 && weekday > 0) {
                    for (let i = 0; i < weekday; i++) {
                        currentWeek.push(null);
                    }
                }

                currentWeek.push(dayData);
            });

            // Push the last week
            if (currentWeek.length > 0) {
                weeks.push(currentWeek);
            }

            // Render all days from all weeks
            weeks.forEach(week => {
                week.forEach(dayData => {
                    if (dayData === null) {
                        // Empty cell for alignment
                        const emptyDiv = document.createElement('div');
                        emptyDiv.className = 'day-cell no-slots';
                        emptyDiv.style.visibility = 'hidden';
                        calendar.appendChild(emptyDiv);
                        return;
                    }

                    const dayDiv = document.createElement('div');
                    const availableCount = dayData.slots.filter(s => s.available).length;
                    const hasAvailable = availableCount > 0;

                    dayDiv.className = 'day-cell' + (hasAvailable ? '' : ' no-slots');

                    const date = dayData.date;
                    const dayName = date.toLocaleDateString('en-US', { weekday: 'short' });
                    const dayNumber = date.getDate();

                    dayDiv.innerHTML = '<div class="day-number">' + dayNumber + '</div>' +
                        '<div class="day-name">' + dayName + '</div>' +
                        '<div class="day-slots-count">' + availableCount + ' available</div>';

                    if (hasAvailable) {
                        dayDiv.onclick = () => selectDay(dayData);
                    }

                    calendar.appendChild(dayDiv);
                });
            });

            updateMonthNavigation();
        }

        function selectDay(dayData) {
            selectedDay = dayData;
            document.getElementById('calendar').style.display = 'none';
            document.getElementById('monthNavigation').style.display = 'none';
            document.getElementById('timeSlotsPanel').classList.add('active');

            const dateStr = dayData.date.toLocaleDateString('en-US', {
                weekday: 'long',
                month: 'long',
                day: 'numeric',
                year: 'numeric'
            });
            document.getElementById('selectedDateTitle').textContent = dateStr;

            renderTimeSlots(dayData.slots);
            document.getElementById('timeSlotsPanel').scrollIntoView({ behavior: 'smooth' });
        }

        function renderTimeSlots(slots) {
            const grid = document.getElementById('timeSlotsGrid');
            grid.innerHTML = '';

            slots.sort((a, b) => new Date(a.slot_time) - new Date(b.slot_time));

            slots.forEach(slot => {
                const timeSlot = document.createElement('div');
                timeSlot.className = 'time-slot' + (slot.available ? '' : ' booked');

                timeSlot.textContent = formatTime(slot.slot_time);

                if (slot.available) {
                    timeSlot.onclick = () => selectTimeSlot(slot);
                }

                grid.appendChild(timeSlot);
            });
        }

        function selectTimeSlot(slot) {
            selectedSlot = slot;
            document.getElementById('selectedSlotDisplay').textContent = formatDateTime(slot.slot_time);
            document.getElementById('timeSlotsPanel').classList.remove('active');
            document.getElementById('bookingForm').classList.add('active');
            document.getElementById('bookingForm').scrollIntoView({ behavior: 'smooth' });
        }

        function backToCalendar() {
            document.getElementById('timeSlotsPanel').classList.remove('active');
            document.getElementById('calendar').style.display = 'grid';
            document.getElementById('monthNavigation').style.display = 'flex';
            selectedDay = null;
        }

        function cancelBooking() {
            selectedSlot = null;
            document.getElementById('bookingForm').classList.remove('active');
            document.getElementById('name').value = '';
            document.getElementById('email').value = '';

            if (selectedDay) {
                document.getElementById('timeSlotsPanel').classList.add('active');
            } else {
                backToCalendar();
            }
        }

        async function confirmBooking() {
            const name = document.getElementById('name').value.trim();
            const email = document.getElementById('email').value.trim();

            if (!name || !email) {
                showMessage('Please fill in all fields', 'error');
                return;
            }

            if (!validateEmail(email)) {
                showMessage('Please enter a valid email address', 'error');
                return;
            }

            try {
                const response = await fetch('/api/bookings', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        slot_time: selectedSlot.slot_time,
                        name: name,
                        email: email
                    })
                });

                if (response.ok) {
                    showMessage('Booking confirmed! You will receive a confirmation email shortly.', 'success');

                    // Mark the slot as unavailable in local data
                    markSlotAsBooked(selectedSlot.slot_time);

                    // Clear form and selections
                    const wasSelectedDay = selectedDay;
                    cancelBooking();

                    // Return to calendar view with updated data
                    if (wasSelectedDay) {
                        backToCalendar();
                    }
                    renderCalendar();
                } else if (response.status === 409) {
                    showMessage('This slot has already been booked. Please select another time.', 'error');
                    loadSlots(); // Reload slots to get fresh data
                } else {
                    const error = await response.text();
                    showMessage('Failed to create booking: ' + error, 'error');
                }
            } catch (error) {
                console.error('Error creating booking:', error);
                showMessage('Failed to create booking. Please try again.', 'error');
            }
        }

        function markSlotAsBooked(slotTime) {
            // Update in allSlots array
            const slot = allSlots.find(s => s.slot_time === slotTime);
            if (slot) {
                slot.available = false;
            }

            // Update in availableMonths structure
            availableMonths.forEach(monthData => {
                monthData.days.forEach(dayData => {
                    const slotInDay = dayData.slots.find(s => s.slot_time === slotTime);
                    if (slotInDay) {
                        slotInDay.available = false;
                    }
                });
            });
        }

        function validateEmail(email) {
            const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            return re.test(email);
        }

        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.className = 'message ' + type + ' active';

            setTimeout(() => {
                messageDiv.classList.remove('active');
            }, 5000);
        }

        // Load slots when page loads
        loadSlots();
    </script>
</body>
</html>
`
	fmt.Fprint(w, html)
}


func AdminHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Admin Panel - Coach Calendar</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }

        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px;
            text-align: center;
        }

        .header h1 {
            font-size: 2.5rem;
            margin-bottom: 10px;
        }

        .header p {
            font-size: 1.1rem;
            opacity: 0.9;
        }

        .nav-link {
            display: inline-block;
            margin-top: 15px;
            padding: 10px 20px;
            background: rgba(255,255,255,0.2);
            color: white;
            text-decoration: none;
            border-radius: 6px;
            transition: all 0.3s;
        }

        .nav-link:hover {
            background: rgba(255,255,255,0.3);
        }

        .content {
            padding: 40px;
        }

        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }

        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 25px;
            border-radius: 12px;
            text-align: center;
        }

        .stat-card h3 {
            font-size: 2rem;
            margin-bottom: 10px;
        }

        .stat-card p {
            font-size: 0.9rem;
            opacity: 0.9;
        }

        .filters {
            display: flex;
            gap: 15px;
            margin-bottom: 30px;
            flex-wrap: wrap;
        }

        .filter-btn {
            padding: 12px 24px;
            border: 2px solid #e0e0e0;
            background: white;
            border-radius: 8px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s;
            font-size: 1rem;
        }

        .filter-btn:hover {
            border-color: #667eea;
            background: #f0f4ff;
        }

        .filter-btn.active {
            background: #667eea;
            color: white;
            border-color: #667eea;
        }

        .slots-container {
            background: #f9f9f9;
            border-radius: 12px;
            padding: 30px;
            max-height: 600px;
            overflow-y: auto;
        }

        .slots-grid {
            display: grid;
            gap: 15px;
        }

        .slot-card {
            background: white;
            padding: 20px;
            border-radius: 10px;
            border: 2px solid #e0e0e0;
            display: grid;
            grid-template-columns: 1fr auto;
            align-items: center;
            gap: 20px;
            transition: all 0.3s;
        }

        .slot-card:hover {
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }

        .slot-card.available {
            border-left: 4px solid #4caf50;
        }

        .slot-card.booked {
            border-left: 4px solid #2196f3;
            background: #e3f2fd;
        }

        .slot-card.blocked {
            border-left: 4px solid #f44336;
            background: #ffebee;
        }

        .slot-info h4 {
            font-size: 1.1rem;
            margin-bottom: 8px;
            color: #333;
        }

        .slot-details {
            display: flex;
            gap: 20px;
            flex-wrap: wrap;
            font-size: 0.9rem;
            color: #666;
        }

        .slot-detail {
            display: flex;
            align-items: center;
            gap: 5px;
        }

        .status-badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 0.8rem;
            font-weight: 600;
            text-transform: uppercase;
        }

        .status-badge.available {
            background: #c8e6c9;
            color: #2e7d32;
        }

        .status-badge.booked {
            background: #bbdefb;
            color: #1565c0;
        }

        .status-badge.blocked {
            background: #ffcdd2;
            color: #c62828;
        }

        .slot-actions {
            display: flex;
            gap: 10px;
        }

        .action-btn {
            padding: 10px 20px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s;
            font-size: 0.9rem;
        }

        .action-btn.block {
            background: #f44336;
            color: white;
        }

        .action-btn.block:hover {
            background: #d32f2f;
        }

        .action-btn.unblock {
            background: #4caf50;
            color: white;
        }

        .action-btn.unblock:hover {
            background: #388e3c;
        }

        .action-btn:disabled {
            background: #ccc;
            cursor: not-allowed;
        }

        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
            font-size: 1.1rem;
        }

        .message {
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
            text-align: center;
            display: none;
        }

        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        .message.active {
            display: block;
        }

        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #999;
        }

        .empty-state h3 {
            font-size: 1.5rem;
            margin-bottom: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Admin Panel</h1>
            <p>Manage bookings and slot availability</p>
            <a href="/" class="nav-link">← Back to Booking Page</a>
        </div>

        <div class="content">
            <div id="message" class="message"></div>

            <div class="stats">
                <div class="stat-card">
                    <h3 id="totalSlots">-</h3>
                    <p>Total Slots</p>
                </div>
                <div class="stat-card">
                    <h3 id="availableSlots">-</h3>
                    <p>Available</p>
                </div>
                <div class="stat-card">
                    <h3 id="bookedSlots">-</h3>
                    <p>Booked</p>
                </div>
                <div class="stat-card">
                    <h3 id="blockedSlots">-</h3>
                    <p>Blocked</p>
                </div>
            </div>

            <div class="filters">
                <button class="filter-btn active" onclick="filterSlots('all')">All Slots</button>
                <button class="filter-btn" onclick="filterSlots('available')">Available</button>
                <button class="filter-btn" onclick="filterSlots('booked')">Booked</button>
                <button class="filter-btn" onclick="filterSlots('blocked')">Blocked</button>
            </div>

            <div id="loading" class="loading">
                Loading slots...
            </div>

            <div id="slotsContainer" class="slots-container" style="display: none;">
                <div id="slotsGrid" class="slots-grid"></div>
            </div>
        </div>
    </div>

    <script>
        let allSlots = [];
        let currentFilter = 'all';

        function formatDateTime(isoString) {
            const date = new Date(isoString);
            return date.toLocaleString('en-US', {
                weekday: 'long',
                month: 'long',
                day: 'numeric',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        }

        async function loadSlots() {
            try {
                const response = await fetch('/api/admin/slots');
                if (!response.ok) {
                    throw new Error('Failed to load slots');
                }

                allSlots = await response.json();
                updateStats();
                renderSlots();

                document.getElementById('loading').style.display = 'none';
                document.getElementById('slotsContainer').style.display = 'block';
            } catch (error) {
                console.error('Error loading slots:', error);
                showMessage('Failed to load slots. Please refresh the page.', 'error');
                document.getElementById('loading').style.display = 'none';
            }
        }

        function updateStats() {
            const stats = {
                total: allSlots.length,
                available: allSlots.filter(s => s.status === 'available').length,
                booked: allSlots.filter(s => s.status === 'booked').length,
                blocked: allSlots.filter(s => s.status === 'blocked').length
            };

            document.getElementById('totalSlots').textContent = stats.total;
            document.getElementById('availableSlots').textContent = stats.available;
            document.getElementById('bookedSlots').textContent = stats.booked;
            document.getElementById('blockedSlots').textContent = stats.blocked;
        }

        function filterSlots(filter) {
            currentFilter = filter;

            // Update button states
            document.querySelectorAll('.filter-btn').forEach(btn => {
                btn.classList.remove('active');
            });
            event.target.classList.add('active');

            renderSlots();
        }

        function renderSlots() {
            const grid = document.getElementById('slotsGrid');
            grid.innerHTML = '';

            const filteredSlots = currentFilter === 'all'
                ? allSlots
                : allSlots.filter(s => s.status === currentFilter);

            if (filteredSlots.length === 0) {
                grid.innerHTML = '<div class="empty-state"><h3>No slots found</h3><p>Try adjusting your filter.</p></div>';
                return;
            }

            filteredSlots.forEach(slot => {
                const slotCard = document.createElement('div');
                slotCard.className = 'slot-card ' + slot.status;

                let detailsHTML = '<span class="status-badge ' + slot.status + '">' + slot.status + '</span>';

                if (slot.status === 'booked' && slot.name && slot.email) {
                    detailsHTML += '<div class="slot-detail">👤 ' + slot.name + '</div>';
                    detailsHTML += '<div class="slot-detail">📧 ' + slot.email + '</div>';
                }

                let actionsHTML = '';
                if (slot.status === 'available') {
                    actionsHTML = '<button class="action-btn block" onclick="blockSlot(\'' + slot.slot_time + '\')">Block</button>';
                } else if (slot.status === 'blocked') {
                    actionsHTML = '<button class="action-btn unblock" onclick="unblockSlot(\'' + slot.slot_time + '\')">Unblock</button>';
                }

                slotCard.innerHTML =
                    '<div class="slot-info">' +
                        '<h4>' + formatDateTime(slot.slot_time) + '</h4>' +
                        '<div class="slot-details">' + detailsHTML + '</div>' +
                    '</div>' +
                    '<div class="slot-actions">' + actionsHTML + '</div>';

                grid.appendChild(slotCard);
            });
        }

        async function blockSlot(slotTime) {
            try {
                const response = await fetch('/api/admin/block', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ slot_time: slotTime })
                });

                if (response.ok) {
                    showMessage('Slot blocked successfully', 'success');
                    await loadSlots();
                } else {
                    const error = await response.text();
                    showMessage('Failed to block slot: ' + error, 'error');
                }
            } catch (error) {
                console.error('Error blocking slot:', error);
                showMessage('Failed to block slot. Please try again.', 'error');
            }
        }

        async function unblockSlot(slotTime) {
            try {
                const response = await fetch('/api/admin/unblock', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ slot_time: slotTime })
                });

                if (response.ok) {
                    showMessage('Slot unblocked successfully', 'success');
                    await loadSlots();
                } else {
                    const error = await response.text();
                    showMessage('Failed to unblock slot: ' + error, 'error');
                }
            } catch (error) {
                console.error('Error unblocking slot:', error);
                showMessage('Failed to unblock slot. Please try again.', 'error');
            }
        }

        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.className = 'message ' + type + ' active';

            setTimeout(() => {
                messageDiv.classList.remove('active');
            }, 5000);
        }

        // Load slots when page loads
        loadSlots();
    </script>
</body>
</html>
`
	fmt.Fprint(w, html)
}


func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
