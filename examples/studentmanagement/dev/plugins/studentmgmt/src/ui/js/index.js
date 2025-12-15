// Initialize method for studentmgmt UI plugin
// This is called by reactapplication when the plugin is loaded

// Import sagas and reducers
import './sagas/StudentSagas';
import './reducers/StudentData';

var module;

function Initialize(appName, ins, mod, settings, def, req) {
    module = this;

    // Access to properties and localization
    module.properties = Application.Properties[ins];
    module.settings = settings;

    console.log('Student Management UI plugin initialized');
}

export {
    Initialize
}
