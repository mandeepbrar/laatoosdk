# Design System

Laatoo provides a set of reusable UI components in the `reactwebcommon` plugin. These components ensure a consistent look and feel across applications.

## Usage

Import components from `reactwebcommon`:

```javascript
import { Card, TextField, ActionButton } from 'reactwebcommon/components/ui';
```

## Components

### Card

Used to display content in a container.

```javascript
<Card title="Student Profile">
    <div>Name: John Doe</div>
    <div>ID: 12345</div>
</Card>
```

### ActionButton

A button for performing actions.

```javascript
<ActionButton 
    label="Register" 
    onClick={() => handleRegister()} 
    primary={true} 
/>
```

### TextField

Input field for text.

```javascript
<TextField 
    label="Course Name" 
    value={courseName} 
    onChange={(e) => setCourseName(e.target.value)} 
/>
```

### Select

Dropdown selection.

```javascript
<Select 
    label="Department" 
    options={[{label: 'CS', value: 'cs'}, {label: 'Math', value: 'math'}]} 
    value={department} 
    onChange={(val) => setDepartment(val)} 
/>
```

### Checkbox

Boolean input.

```javascript
<Checkbox 
    label="Active" 
    checked={isActive} 
    onChange={(val) => setIsActive(val)} 
/>
```

### Switch

Toggle switch.

```javascript
<Switch 
    label="Notifications" 
    checked={notificationsEnabled} 
    onChange={(val) => setNotificationsEnabled(val)} 
/>
```

### Tabs

Tabbed interface.

```javascript
<Tabs 
    tabs={[{label: 'Details', id: 'details'}, {label: 'Courses', id: 'courses'}]} 
    activeTab={activeTab} 
    onTabChange={(id) => setActiveTab(id)} 
/>
```

### Navbar

Navigation bar.

```javascript
<Navbar 
    title="Student Management" 
    items={[{label: 'Home', link: '/'}, {label: 'Students', link: '/students'}]} 
/>
```

### Avatar

User avatar.

```javascript
<Avatar 
    src="/path/to/image.jpg" 
    name="John Doe" 
/>
```

### Icon

Display icons.

```javascript
<Icon name="user" />
```

## Example: Student Registration Form

```javascript
import React, { useState } from 'react';
import { Card, TextField, Select, ActionButton } from 'reactwebcommon';

export default function StudentRegistration() {
    const [name, setName] = useState('');
    const [dept, setDept] = useState('');

    const handleSubmit = () => {
        // Submit logic
    };

    return (
        <Card title="Register Student">
            <TextField 
                label="Full Name" 
                value={name} 
                onChange={(e) => setName(e.target.value)} 
            />
            <Select 
                label="Department" 
                options={[
                    {label: 'Computer Science', value: 'CS'},
                    {label: 'Mathematics', value: 'MATH'}
                ]} 
                value={dept} 
                onChange={(val) => setDept(val)} 
            />
            <ActionButton 
                label="Submit" 
                onClick={handleSubmit} 
                primary={true} 
            />
        </Card>
    );
}
```
