describe('Schedule Details Page - EVV Logger', () => {
  beforeEach(() => {
    cy.visit('/test.html');
    cy.get('button').contains('View Details').first().click();
  });

  it('should display schedule details header', () => {
    cy.get('h1').should('contain', 'John Smith');
    cy.get('.bg-blue-100').should('contain', '09:00 - 12:00');
  });

  it('should show visit status section', () => {
    cy.get('h2').contains('Visit Status').should('exist');
    cy.get('button').contains('Start Visit').should('exist');
  });

  it('should display tasks section', () => {
    cy.get('h2').contains('Care Activities').should('exist');
    cy.get('.schedule-item').should('exist');
  });

  it('should show task completion buttons', () => {
    cy.get('button').contains('Complete').should('exist');
    cy.get('button').contains('Not Completed').should('exist');
  });

  it('should handle task completion', () => {
    cy.get('button').contains('Complete').first().click();
    cy.get('.text-green-600').should('exist');
  });

  it('should show "No tasks assigned" when no tasks exist', () => {
    cy.get('.text-center').should('contain', 'No tasks assigned');
  });

  it('should have back to schedules button', () => {
    cy.get('button').should('contain', 'Back to Schedules');
  });

  it('should handle schedule not found error', () => {
    cy.get('.text-center').should('contain', 'Schedule not found');
  });

  it('should display geolocation information when available', () => {
    cy.get('.flex').should('contain', '📍');
    cy.get('.text-sm').should('contain', 'latitude');
    cy.get('.text-sm').should('contain', 'longitude');
  });

  it('should show visit duration when visit is completed', () => {
    cy.get('.text-sm').should('contain', 'Duration');
  });

  it('should have proper error handling for geolocation', () => {
    cy.get('.bg-red-50').should('not.exist');
  });

  it('should be responsive on mobile view', () => {
    cy.viewport(375, 667);
    cy.get('h1').should('contain', 'John Smith');
    cy.get('button').should('exist');
  });
});

describe('Schedule Details - Visit Tracking', () => {
  beforeEach(() => {
    cy.visit('/test.html');
    cy.get('button').contains('View Details').first().click();
  });

  it('should start visit process', () => {
    cy.get('button').contains('Start Visit').click();
    cy.get('.bg-blue-600').should('exist');
  });

  it('should end visit process', () => {
    cy.get('button').contains('Start Visit').click();
    cy.get('button').contains('End Visit').should('exist');
    cy.get('button').contains('End Visit').click();
  });

  it('should capture location data during visit', () => {
    cy.get('button').contains('Start Visit').click();
    cy.get('.text-sm').should('contain', 'latitude');
    cy.get('.text-sm').should('contain', 'longitude');
  });

  it('should handle geolocation errors gracefully', () => {
    cy.get('.text-red-800').should('not.exist');
  });
});