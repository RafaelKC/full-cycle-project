import { Module } from '@nestjs/common';
import { RoutesService } from './routes.service';
import { RoutesController } from './routes.controller';
import { PrismaModule } from '../prisma/prisma.module';
import { MapsModule } from '../maps/maps.module';
import { RoutesDriverService } from './routes-driver/routes-driver.service';

@Module({
  imports: [PrismaModule, MapsModule],
  controllers: [RoutesController],
  providers: [RoutesService, RoutesDriverService],
})
export class RoutesModule {}
