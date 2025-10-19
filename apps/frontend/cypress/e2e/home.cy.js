describe('Home Page - EVV Logger', () => {
  beforeEach(() => {
    cy.visit('/test.html');
  });

  it('should display the EVV Logger title', () => {
    cy.get('h1').should('contain', 'EVV Logger');
  });

  it('should show stats dashboard with all cards', () => {
    cy.get('.stats-card').should('have.length', 4);
    cy.get('.stats-card').first().should('contain', 'Total Schedules');
    cy.get('.stats-card').eq(1).should('contain', 'Upcoming');
    cy.get('.stats-card').eq(2).should('contain', 'Completed');
    cy.get('.stats-card').last().should('contain', 'Missed');
  });

  it('should display schedule list with sample data', () => {
    cy.get('.schedule-item').should('have.length.greaterThan', 0);
    cy.get('.schedule-item').first().should('contain', 'John Smith');
    cy.get('.schedule-item').first().should('contain', '09:00 - 12:00');
  });

  it('should have status filtering dropdown', () => {
    cy.get('select').should('exist');
    cy.get('select').should('contain', 'All Statuses');
    cy.get('select').should('contain', 'Upcoming');
    cy.get('select').should('contain', 'Completed');
    cy.get('select').should('contain', 'Missed');
  });

  it('should filter schedules by status', () => {
    // Select "Completed" status
    cy.get('select').select('Completed');

    // Should only show completed schedules
    cy.get('.schedule-item').each(($el) => {
      cy.wrap($el).should('contain', 'Completed');
    });
  });

  it('should have view details buttons', () => {
    cy.get('button').contains('View Details').should('exist');
    cy.get('button').contains('View Details').first().click();
  });

  it('should handle empty schedule state', () => {
    // This test would need to mock the API to return empty data
    cy.get('.schedule-item').should('have.length', 0);
  });

  it('should be responsive on mobile view', () => {
    cy.viewport(375, 667); // iPhone 6/7/8
    cy.get('.stats-card').should('have.length', 4);
    cy.get('.schedule-item').should('exist');
  });

  it('should have proper accessibility attributes', () => {
    cy.get('h1').should('have.attr', 'role');
    cy.get('button').should('have.attr', 'role', 'button');
  });
});